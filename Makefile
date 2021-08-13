XGOARCH := amd64
XGOOS := linux
XBIN := $(XGOOS)_$(XGOARCH)/atb

all: test vet tools install

fmt:
	go fmt ./...

test:
	go test ./...

vet:
	go vet ./...

# https://github.com/golang/go/issues/25922
# https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
tools:
	go generate -tags tools ./...

install:
	go install ./...

xinstall:
	env GOOS=$(XGOOS) GOARCH=$(XGOARCH) go install ./...

publish:
ifndef DEST_PATH
	$(error DEST_PATH must be set when publishing)
endif
	rsync -az $(GOPATH)/bin/$(XBIN) $(DEST_PATH)/$(XBIN)
	@sha256sum $(GOPATH)/bin/$(XBIN)
