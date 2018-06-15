clean:
	rm -Rfv bin
	mkdir bin

build: clean
	go build -o bin/gitlab-exporter *.go

build-all: clean
	GOOS="linux"   GOARCH="amd64"       go build -o bin/gitlab-exporter__linux-amd64 *.go
	GOOS="linux"   GOARCH="arm" GOARM=6 go build -o bin/gitlab-exporter__linux-armv6 *.go
	GOOS="linux"   GOARCH="arm" GOARM=7 go build -o bin/gitlab-exporter__linux-armv7 *.go
	GOOS="linux"   GOARCH="arm"         go build -o bin/gitlab-exporter__linux-arm   *.go
	GOOS="darwin"  GOARCH="amd64"       go build -o bin/gitlab-exporter__macos-amd64 *.go
	GOOS="windows" GOARCH="amd64"       go build -o bin/gitlab-exporter__win-amd64   *.go
