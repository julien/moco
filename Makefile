SHELL := /bin/bash
COVFILE := cover.out

fmt:
	gofmt -w *.go

vet: fmt
	go vet

test: fmt
	go test -coverprofile=$(COVFILE) -p 1 -race

cover: test
	go tool cover -html=$(COVFILE)
