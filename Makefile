export GOPROXY=https://goproxy.io

default: build

build: export GO111MODULE=on

build:
	GOOS=linux go build -o bin/secret-game cmd/main.go
	