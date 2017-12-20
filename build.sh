#!/bin/bash

rm -f utils/assets.go && go generate ./... && go build -o bin/serverbutler