REPOTAG=embano1/simplelog:1.0

default: build

build:
	GOOS=linux CGO_ENABLED=0 go build -a -installsuffix "static" -o web cmd/web/*.go
	GOOS=linux CGO_ENABLED=0 go build -a -installsuffix "static" -o tailf cmd/tailf/*.go

image: 
	docker build -t ${REPOTAG} .

push:
	docker push ${REPOTAG}

clean:
	rm web tailf

all: build image clean push
	