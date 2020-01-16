SHELL := /bin/bash

.PHONY: nomad-dtree

clean: 
	rm -rf ./dist/nomad-dtree

build:
	GOOS=linux GARCH=amd64 CGO_ENABLED=0 \
		go build -o ./dist/nomad-dtree -a -installsuffix cgo .

test:
	go test -v ./...
