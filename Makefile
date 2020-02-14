.PHONY: test build

TAG := $(shell git rev-parse --short HEAD)

REGISTRIES?=registry.qtt6.cn
REPOSITORY?=paas-dev
APP=slabinfo
image=$(shell cat release)

test:
	go test -mod=vendor ./...

build:
	GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -o deploy/bin/$(APP) main.go

image: build
	@docker build -f deploy/Dockerfile deploy -t $(REGISTRIES)/$(REPOSITORY)/$(APP)
	@docker push $(REGISTRIES)/$(REPOSITORY)/$(APP)
	@echo "$(REGISTRIES)/$(REPOSITORY)/$(APP)"

