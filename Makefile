# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = lazycopy
GOARCH = amd64

VERSION=1.0.2
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
#GITHUB_USERNAME=amimof
#BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: clean fmt linux darwin windows

linux: 
	go get ./... ; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} .

rpi: 
	go get ./... ; \
	GOOS=linux GOARCH=arm go build ${LDFLAGS} -o ${BINARY}-linux-arm .

darwin:
	go get ./... ; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} .

windows:
	go get ./... ; \
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe .

test:
	go get ./... ; \
	go test

fmt:
	go get ./... ; \
	go fmt $$(go list ./... | grep -v /vendor/)

clean:
	-rm -f ${BINARY}-*

.PHONY: linux rpi darwin windows test fmt clean