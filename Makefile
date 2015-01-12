SHELL := /bin/bash
COVFILE := cover.out

fmt:
	go fmt

lint: fmt
	golint *.go

vet: lint
	go vet

test: lint
	go test -coverprofile=$(COVFILE)

cover: test
	go tool cover -html=$(COVFILE)
