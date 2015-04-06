#!/usr/bin/env bash

set -e

function lint() {
    go get -v github.com/golang/lint/golint
    golint ./...
}

function vet() {
    local -r flags="$1"
    go get -v golang.org/x/tools/cmd/vet
    go tool vet $flags $PWD
}


case "$TRAVIS_GO_VERSION" in
    1.1*)
        printf "go ${TRAVIS_GO_VERSION} doesn't support lint or vet\n"
        ;;
    *)
        lint
        vet
        ;;
esac
