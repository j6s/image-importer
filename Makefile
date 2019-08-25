clean:
	rm -Rfv bin
	mkdir bin

build: clean
	gcc ./src/main.c \
		-o bin/image-importer \
		-lexif

# TODO: Cross compilation
#build-all: clean
#	GOOS="linux"   GOARCH="amd64"       go build -o bin/image-importer__linux-amd64 cmd/image-importer/*.go
#	GOOS="linux"   GOARCH="arm" GOARM=6 go build -o bin/image-importer__linux-armv6 cmd/image-importer/*.go
#	GOOS="linux"   GOARCH="arm" GOARM=7 go build -o bin/image-importer__linux-armv7 cmd/image-importer/*.go
#	GOOS="linux"   GOARCH="arm"         go build -o bin/image-importer__linux-arm   cmd/image-importer/*.go
#	GOOS="darwin"  GOARCH="amd64"       go build -o bin/image-importer__macos-amd64 cmd/image-importer/*.go
#	GOOS="windows" GOARCH="amd64" go build -o bin/image-importer__win-amd64 cmd/image-importer/*.go