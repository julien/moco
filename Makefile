SHELL := /bin/bash
COVFILE := cover.out

fmt:
	gofmt -w *.go

vet: fmt
	go vet

test: fmt
	go test -coverprofile=$(COVFILE)

cover: test
	go tool cover -html=$(COVFILE)
