all: deps test vet lint install

fmt:
	go fmt ./...

test:
	go test ./...

vet:
	go vet ./...

lint:
	golint 2> /dev/null; if [ $$? -eq 127 ]; then \
		go get -v github.com/golang/lint/golint; \
	fi
	golint ./...

deps:
	go get -d -v ./...

install:
	go install
