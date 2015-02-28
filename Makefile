NAME=atbapi

all: deps install

deps:
	go get -d -v

install:
	go install

docker-build:
	docker run --rm -v $(PWD):/usr/src/$(NAME) -w /usr/src/$(NAME) \
		golang:latest /bin/sh -c 'go get -d -v && go build -v'
