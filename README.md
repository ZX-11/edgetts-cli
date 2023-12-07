# edgetts-cli

利用Edge提供的tts接口进行语音合成。支持任意长度的文本，支持任意输出音频格式（配合ffmpeg），支持并行抓取，支持断点续传。

本项目基于[https://github.com/jing332/tts-server-go/](https://github.com/jing332/tts-server-go/)开发。

Windows用户可以直接下载可执行文件使用，下载后将压缩包解压，并将可执行文件加入到PATH即可：

> <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16c0 1.1.9 2 2 2h12a2 2 0 0 0 2-2V8l-6-6z"/><path d="M14 3v5h5M16 13H8M16 17H8M10 9H8"/></svg><a href="https://github.com/ZX-11/edgetts-cli/releases/download/0.1/edgetts.zip"><b>edgetts.zip</b></a>

Linux用户可以通过脚本构建（需要安装go和ffmpeg）:

```bash
# for debian users
apt install golang ffmpeg

# for archlinux users
pacman -S go ffmpeg

curl https://fastly.jsdelivr.net/gh/ZX-11/edgetts-cli@main/build.sh | sh
```

选项说明：
- -i 需要合成的文本文件，须为UTF-8编码（不能含有BOM），且每段文本不能超出1000字
- -o 输出的音频文件，默认输出opus编码的音频文件，拓展名须为ogg/opus/webm
- -rate 语音速度，可以为x-slow/slow/medium/fast/x-fast或一个相对值，默认为medium
- -voice 合成声音，默认为zh-CN-XiaoxiaoNeural，也有其他声音可选，详见`edgetts -h`
- -parallel 并行连接，默认为1，最大为8，提高并行连接数可适当加快合成速度，但需要避免并行连接过大导致滥用
- -convert 转换为其他格式（即可在-o选项中使用其他格式如mp3/m4a/amr等），需要外部ffmpeg程序（linux用户需要安装ffmpeg，windows用户可以从官网下载并配置，也可使用everything等工具搜索计算机中现有ffmpeg.exe可执行文件并加入到PATH，很多常用影视软件如优酷客户端、格式工厂均包含ffmpeg.exe）

使用示例：
- 合成text.txt为语音，输出为out.ogg：`edgetts -i text.txt -o out.ogg`
- 使用4个连接合成text.txt为语音，输出为out.ogg：`edgetts -i text.txt -o out.ogg -parallel 4`
- 指定声音合成：`edgetts -i text.txt -voice zh-CN-YunyangNeural -o out.ogg`
- 调整声音速度：`edgetts -i text.txt -rate fast -o out.ogg`
- 生成MP3格式音频：`edgetts -i text.txt -convert -o out.mp3`
- 组合各选项使用：`edgetts -i text.txt -voice zh-CN-XiaoshuangNeural -rate slow -convert -o out.m4a -parallel 4`

常见报错处理：
- `ffmpeg not found`：请确认edgetts可执行文件同目录下存在ffmpeg-min可执行文件
- `external ffmpeg not found`：请确认ffmpeg可执行文件位于PATH中
- `Invalid utf-8 sequence`：请将文件从GBK编码转换为UTF-8编码，最简单的方式是用记事本打开后另存为，编码选择`UTF-8`
- `Too long for line xxx`：某段文字过长，需要适当换行分割

**请合理使用该工具，切勿滥用导致接口失效，日常还是建议尽量在Edge浏览器中使用该功能。**
