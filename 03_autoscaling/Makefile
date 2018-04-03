REPO=embano1/pubsub_autoscaler
VERSION=1.5-vmworld

build:
	for bin in sender receiver autoscaler ; \
		do GOOS=linux go build -a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo -o out/$${bin} cmd/$${bin}/*.go ; \
	done

all: image

image: build
	docker build -t ${REPO}:${VERSION} .

clean: 
# wont remove the docker image
	rm out/*
	
