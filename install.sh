type wget && \
type unzip && \
type go && \
wget https://github.com/jing332/tts-server-go/archive/refs/heads/master.zip --no-check-certificate && \
unzip master.zip && \
rm master.zip && \
cd tts-server-go-master  && \
mkdir edgetts-cli && \
cd edgetts-cli && \
wget https://github.com/ZX-11/edgetts-cli/raw/main/edgetts.go --no-check-certificate && \
go build -ldflags="-s -w" edgetts.go && \
echo "All Done. Build file can be found in ./tts-server-go-master/edgetts-cli"