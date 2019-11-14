.PHONY: all test

all: check test

check: goimports govet

goimports:
	goimports -d .

govet:
	go vet ./...

test:
	go test -v ./...

build:
	go build -o gfdupes .

install:
	go install

clean:
	go clean .

bootstrap:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
