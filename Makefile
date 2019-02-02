.PHONY: build
build: generate
	CGO_ENABLED=0 GOOS=linux go build -o bin/serverbutler

.PHONY: clean
clean:
	rm -f utils/assets.go

.PHONY: generate
generate: clean
	go generate ./...

.PHONY: docker-build
docker: build
	docker build -t serverbutler .

.PHONY: run
run: clean generate
	go build -o bin/serverbutler
	bin/serverbutler
