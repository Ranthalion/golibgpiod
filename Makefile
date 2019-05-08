VERSION=$(shell git describe --always --tags)

.PHONY: go
go:
	go build -o bin/gogpiodetect -ldflags "-s -w -X main.Version=$(VERSION)"  ./cmd/gogpiodetect
	go build -o bin/gogpioinfo -ldflags "-s -w -X main.Version=$(VERSION)"  ./cmd/gogpioinfo

.PHONY: clean
clean:
	rm -rf bin/*
