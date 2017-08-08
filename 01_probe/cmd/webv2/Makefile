REPO = embano1/simpleweb
VERSION = 2.0
DATE := $(shell date -u +%D-%H%M%S)
REV := $(shell git rev-parse --short HEAD)

default: build

build:
	GOOS=linux CGO_ENABLED=0 go build -tags netgo -ldflags "-X main.version=${VERSION} -X main.build=${REV}@${DATE}" -o webv2

image: 
	docker build -t ${REPO}:${VERSION} .

push:
	docker push ${REPO}:${VERSION}

clean:
	rm webv2

all: build image clean push
	
