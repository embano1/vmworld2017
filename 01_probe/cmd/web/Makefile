REPO = embano1/simpleweb
VERSION = 1.0
DATE := $(shell date -u +%D-%H%M%S)
REV := $(shell git rev-parse --short HEAD)

default: build

build:
	GOOS=linux CGO_ENABLED=0 go build -tags netgo -ldflags "-X main.version=${VERSION} -X main.build=${REV}@${DATE}" .

image: 
	docker build -t ${REPO}:${VERSION} .

push:
	docker push ${REPO}:${VERSION}

clean:
	rm web

all: build image clean push
	