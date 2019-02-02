.PHONY: clean
all: clean generate build

.PHONY: clean
clean:
	rm -f utils/assets.go

.PHONY: generate
generate:
	go generate ./...

.PHONY: run
run: generate
	go build -o bin/serverbutler
	bin/serverbutler

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/serverbutler

.PHONY: docker-build
docker: build
	docker build -t serverbutler .