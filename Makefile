build:
	go build -o office-neighbors -ldflags "-X main.BuildTime=`date +%Y-%m-%d:%H:%M:%S` -X main.GitHash=`git rev-parse --short HEAD`"

all: build