type wget && \
type unzip && \
type go && (\
wget --timeout 10 --tries 1 https://github.com/jing332/tts-server-go/archive/refs/heads/master.zip || \
wget --timeout 10 --tries 1 https://gh.con.sh/https://github.com/jing332/tts-server-go/archive/refs/heads/master.zip || \
wget --timeout 10 --tries 1 https://archive.fastgit.org/jing332/tts-server-go/archive/refs/heads/master.zip || \
wget --timeout 10 --tries 1 https://ghproxy.com/https://github.com/jing332/tts-server-go/archive/refs/heads/master.zip ) && \
unzip master.zip && \
rm master.zip && \
cd tts-server-go-master  && \
mkdir edgetts-cli && \
cd edgetts-cli && (\
wget --timeout 10 --tries 1 https://github.com/ZX-11/edgetts-cli/raw/main/edgetts.go || \
wget --timeout 10 --tries 1 https://cdn.staticaly.com/gh/ZX-11/edgetts-cli/main/edgetts.go || \
wget --timeout 10 --tries 1 https://testingcf.jsdelivr.net/gh/ZX-11/edgetts-cli@main/edgetts.go || \
wget --timeout 10 --tries 1 https://fastly.jsdelivr.net/gh/ZX-11/edgetts-cli@main/edgetts.go || \
wget --timeout 10 --tries 1 https://gcore.jsdelivr.net/gh/ZX-11/edgetts-cli@main/edgetts.go || \
wget --timeout 10 --tries 1 https://raw.fastgit.org/ZX-11/edgetts-cli/main/edgetts.go || \
wget --timeout 10 --tries 1 https://ghproxy.com/https://github.com/ZX-11/edgetts-cli/raw/main/edgetts.go ) && \
go build -ldflags="-s -w" edgetts.go && \
echo "Done. Build file can be found in ./tts-server-go-master/edgetts-cli" && \
( type ffmpeg || echo "Warning: You also need to install ffmpeg!" )