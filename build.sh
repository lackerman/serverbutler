#!/bin/bash

rm -f utils/assets.go && go generate ./...

CGO_ENABLED=0 GOOS=linux go build -o bin/serverbutler

docker build -t serverbutler .