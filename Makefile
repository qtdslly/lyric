export CURRENT_PATH=$(shell pwd)
export PATH:=$(GOROOT)/bin:$(GOPATH)/bin:$(PATH)
export GOBIN=$(GOROOT)/bin
export GOOS=linux
export GOARCH=amd64

BUILD_TIME=$(shell date "+%Y-%m-%d %H:%M:%S")

time = $(shell date "+%Y-%m-%d %H:%M:%S")

.PHONY: docker motd

all: clean docker motd
	@echo ""
	@echo $(time)
	@echo ""

docker:
	go build -o ./bin/$@ -ldflags="-s -w" ./$@/main.go
	chmod 777 ./bin/$@
	upx -9 ./bin/$@

motd:
	go build -o ./bin/$@ -ldflags="-s -w" ./$@/main.go
	chmod 777 ./bin/$@
	upx -9 ./bin/$@

clean:
	@rm -rf ./bin
	@go clean -cache