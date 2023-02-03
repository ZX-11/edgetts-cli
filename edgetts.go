package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/jing332/tts-server-go/tts/edge"
)

var (
	input  = flag.String("i", "", "Path to the .txt file (UTF-8 encoding)")
	output = flag.String("o", "out.ogg", "Path to the output file (default 48kbps opus audio, only ogg/opus/webm are supported without -convert)")
	voice  = flag.String("voice", "zh-CN-XiaoxiaoNeural", `One of:
	en-US-AriaNeural
	en-US-JennyNeural
	en-US-GuyNeura
	en-US-SaraNeural
	ja-JP-NanamiNeural
	pt-BR-FranciscaNeural
	zh-CN-XiaoxiaoNeural
	zh-CN-YunyangNeural
	zh-CN-YunyeNeural
	zh-CN-YunxiNeural
	zh-CN-XiaohanNeural
	zh-CN-XiaomoNeural
	zh-CN-XiaoxuanNeural
	zh-CN-XiaoruiNeural
	zh-CN-XiaoshuangNeural
	... (other voice supported by edge TTS)
`)
	rate = flag.String("rate", "1", `One of:
	x-slow
	slow
	medium
	fast
	x-fast
	a rate number > 0 (meduim = 1)
	a delta number (+0.5, -0.2, ...)
`)
	parallel = flag.Uint("parallel", 1, "Max download threads (Max 8)")
	convert  = flag.Bool("convert", false, "Output with other formats like mp3/m4a/amr..., external ffmpeg is required")
)

var (
	signal          = struct{}{}
	useLocalFFmpeg  = false
	modified        = true
	localFFmpegPath string
	wg              sync.WaitGroup
	tasks           chan Task
	lines           int
)

type Task struct {
	line         int
	text, output string
}

func main() {
	flag.Parse()

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		ex, _ := os.Executable()
		localFFmpegPath = filepath.Dir(ex) + If(runtime.GOOS == "windows", "/ffmpeg-min.exe", "/ffmpeg-min")
		if _, err := os.Stat(localFFmpegPath); err == nil {
			useLocalFFmpeg = true
			if *convert {
				println("external ffmpeg not found")
				os.Exit(1)
			}
		} else {
			println("ffmpeg not found")
			os.Exit(1)
		}
	}

	buf, err := ioutil.ReadFile(*input)
	if err != nil {
		if input == nil || *input == "" {
			println("Use -h to get usage")
		} else {
			println(err.Error())
		}
		os.Exit(1)
	}

	if *parallel == 0 || *parallel > 8 {
		println("Parallel should between 1-8")
		os.Exit(1)
	}

	if !utf8.Valid(buf) {
		println("Invalid utf-8 sequence")
		os.Exit(1)
	}

	paras := strings.Split(str(buf), If(strings.Contains(str(buf), "\r\n"), "\r\n", "\n"))

	lines = len(paras)

	if *parallel > uint(lines) {
		*parallel = uint(lines)
	}

	for i, para := range paras {
		if len(para) > 3000 {
			println("Too long for line", i+1)
			os.Exit(1)
		}
	}

	partDir := filepath.Base(*input) + ".parts"

	if err := os.MkdirAll(partDir, os.ModePerm); err != nil {
		panic(err)
	}

	configPath := partDir + "/config"

	stat, _ := os.Stat(*input)
	newConfig := fmt.Sprintf("last-modify = %v\nvoice = %v\nrate = %v\n", stat.ModTime().Unix(), *voice, *rate)

	if config, err := os.ReadFile(configPath); err == nil {
		if newConfig == str(config) {
			modified = false
		}
	}

	if modified {
		os.WriteFile(configPath, bytes(newConfig), 0666)
	}

	parts, err := os.Create(partDir + "/index")
	if err != nil {
		panic(err)
	}

	tasks = make(chan Task, *parallel)

	for i := uint(0); i < *parallel; i++ {
		go worker()
	}

	for i, para := range paras {
		line := i + 1
		if len(para) == 0 || len(strings.TrimSpace(para)) == 0 {
			fmt.Printf("Finished: %v/%v (empty)\n", line, lines)
			continue
		}

		fmt.Fprintf(parts, "file '%v'\n", fmt.Sprintf("%v.webm", line))
		savePath := fmt.Sprintf("%v/%v.webm", partDir, line)

		if !modified {
			if _, err := os.Stat(savePath); err == nil {
				fmt.Printf("Finished: %v/%v (skipped)\n", line, lines)
				continue
			}
		}

		wg.Add(1)
		tasks <- Task{line, para, savePath}
		if *parallel == 1 {
			wg.Wait()
		}
	}
	parts.Close()
	wg.Wait()
	close(tasks)

	var cmd *exec.Cmd
	if *convert {
		cmd = exec.Command(If(useLocalFFmpeg, localFFmpegPath, "ffmpeg"), "-y", "-f", "concat", "-i", partDir+"/index", *output)
	} else {
		cmd = exec.Command(If(useLocalFFmpeg, localFFmpegPath, "ffmpeg"), "-y", "-f", "concat", "-i", partDir+"/index", "-c", "copy", *output)
	}
	output, _ := cmd.CombinedOutput()
	fmt.Println(str(output))
}

func worker() {
	tts := &edge.TTS{}
	tts.NewConn()
	for {
		task, ok := <-tasks
		if !ok {
			return
		}
		ssml := `<speak xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="http://www.w3.org/2001/mstts" xmlns:emo="http://www.w3.org/2009/10/emotionml" version="1.0" xml:lang="en-US"><voice name="` + *voice + `"><prosody rate="` + *rate + `" pitch="+0Hz">` + task.text + `</prosody></voice></speak>`
		audioData, err := tts.GetAudio(ssml, "webm-24khz-16bit-mono-opus")
		for err != nil {
			fmt.Printf("Error: %v Retrying...\n", err)
			time.Sleep(time.Second * 3)
			audioData, err = tts.GetAudio(ssml, "webm-24khz-16bit-mono-opus")
		}
		os.WriteFile(task.output, audioData, 0666)
		fmt.Printf("Finished: %v/%v\n", task.line, lines)
		wg.Done()
	}
}

func str(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func bytes(s string) []byte {
	header := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: header.Data,
		Len:  header.Len,
		Cap:  header.Len,
	}))
}

func If[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}
	return falseVal
}
