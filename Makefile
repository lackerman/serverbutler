PROJECTNAME=$(shell basename "$(PWD)")

.PHONY: all
all: clean build

.PHONY: clean
clean:
	rm -rf bin/*

build: main.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) CGO_ENABLED=0 go build -o bin/$(PROJECTNAME)

build-arm: main.go
	@$(MAKE) build GOOS=linux GOARCH=arm GOARM=7

.PHONY: docker
docker:
	docker build -t $(PROJECTNAME) .

.PHONY: run
run: build
	bin/serverbutler

.PHONY: copy
copy: check-server clean docker
	docker save --output $(PROJECTNAME).tar $(PROJECTNAME):latest
	scp $(PROJECTNAME).tar $(SERVER_NAME):.
	rm $(PROJECTNAME).tar

check-server:
ifndef SERVER_NAME
	$(error SERVER_NAME is undefined)
endif
