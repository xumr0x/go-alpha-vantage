.PHONY: build

build:
	go build -o bin/av-crypto cmd/crypto/*.go
	go build -o bin/av-stock cmd/tester/*.go
