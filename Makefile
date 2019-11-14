XGOARCH := amd64
XGOOS := linux
XBIN := $(XGOOS)_$(XGOARCH)/atbapi

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

xinstall:
	env GOOS=$(XGOOS) GOARCH=$(XGOARCH) go install ./...

publish:
ifndef DEST_PATH
	$(error DEST_PATH must be set when publishing)
endif
	rsync -a $(GOPATH)/bin/$(XBIN) $(DEST_PATH)/$(XBIN)
	@sha256sum $(GOPATH)/bin/$(XBIN)
