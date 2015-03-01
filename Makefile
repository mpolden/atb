NAME=atbapi

all: deps test lint install

deps:
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint
	go get -d -v

fmt:
	go fmt ./...

lint:
	go vet ./...
	golint ./...

test:
	go test ./...

install:
	go install

docker-build:
	docker run --rm -v $(PWD):/usr/src/$(NAME) -w /usr/src/$(NAME) \
		golang:latest /bin/sh -c 'go get -d -v && go build -v'
