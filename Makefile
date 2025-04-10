GO_VERSION=$(shell go version)
BUILD_TIME=$(shell date +%F=%Z/%T)
VERSION=$(shell cat ./VERSION.txt)
LDFLAGS=-ldflags "-X main.GoVersion=${GO_VESION} -X main.BuildTime'=${BUILD_TIME} -X main.Version=${VERSION}"

release=no
svn=yes

BUILD_PATH=./src

DEBUG_NAME=bsyslog_d
RELEASE_NAME=bsyslog

PROJECT_NAME=${RELEASE_NAME}
MAIN_PATH=./core/main.go

DEMO_NAME=syslog_server
DEMO_PATH=./demo/server.go

ifeq ($(svn), yes)
SVN_VER=$(shell svn info ./|grep 'Last Changed Rev'|cut -c 19-)
endif

.PHONY: build
build:
	@go get -u ./...
	@echo "gofmt..."
	@gofmt -l -w ./
	@go mod tidy
	@echo "comile..."
	@go build ${LDFLAGS} -a -o ${BUILD_PATH}/${PROJECT_NAME} ${MAIN_PATH}
	@echo "build ${BUILD_PATH}/${PROJECT_NAME}"
	@go build -a -o ${BUILD_PATH}/${DEMO_NAME} ${DEMO_PATH} 
	@echo "build ${BUILD_PATH}/${DEMO_NAME}"
clean:
	@GOBIN=$(GOBIN) go clean
	@rm -rf ${BUILD_PATH}