.PHONY: build
build:
	GOOS=linux CGO=0 go build -o bin/app

.PHONY: clean
clean:

.PHONY: docker-build
docker-build:
	docker build -t serverbutler .
