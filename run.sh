#!/bin/bash

rm -f utils/assets.go \
	&& go generate ./... \
	&& config_dir=./tmp site_prefix=serverbutler go run main.go
