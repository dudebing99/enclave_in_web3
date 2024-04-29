.DEFAULT_GOAL := build-all

export GO111MODULE=on
export DATE=$(shell date '+%Y%m%d%H%M')
#export GOPROXY=https://goproxy.io

enclave-client:
	go build -o bin/enclave-client-$(DATE) ./cmd/client/main.go

enclave-keeper:
	go build -o bin/enclave-keeper-$(DATE) ./cmd/keeper/main.go

all: go-deps enclave-client enclave-keeper

clean:
	@rm -rf bin/

go-deps:
	@mkdir -p bin/
	@cp -r conf bin/
