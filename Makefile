# Set the shell to bash always
SHELL := /bin/bash

build-doc:
	docker build . -f deploy/doc.Dockerfile -t crdsdev/doc:latest

build-gitter:
	docker build . -f deploy/gitter.Dockerfile -t crdsdev/doc-gitter:latest

run-doc:
	go run -v ./cmd/doc

run-gitter:
	go run -v ./cmd/gitter

test:
	go test ./...

.PHONY: build-doc build-gitter