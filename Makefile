.PHONY: all test

all: check test build

check: goimports govet

goimports:
	@echo go imports...
	@goimports -d .

govet:
	@echo go vet...
	@go tool vet .

test:
	@go test -v .

build:
	@echo build gfdupes
	@go build github.com/bruceadowns/gfdupes

clean:
	@go clean github.com/bruceadowns/gfdupes
