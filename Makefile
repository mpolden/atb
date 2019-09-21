all: test vet lint install

fmt:
	go fmt ./...

test:
	go test ./...

vet:
	go vet ./...

lint:
	cd tools && \
		go list -tags tools -f '{{range $$i := .Imports}}{{printf "%s\n" $$i}}{{end}}' | xargs go install
	golint ./...

install:
	go install ./...
