COMMIT := $(shell git describe --dirty --always)
TAG := $(shell git describe --tags $(git rev-list --tags --max-count=1))

LDFLAGS := "-s -w -X main.GitCommit=$(COMMIT)"

.PHONY: build
.DEFAULT_GOAL := build

generate:
	go generate ./...

build: generate;
	CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o dynratelimiter main.go
	
lint:
	golangci-lint run

packandpush:
	docker buildx create --use
	docker buildx build --push --platform linux/amd64,linux/i386,linux/arm64,linux/arm/v6,linux/arm/v7 -t $(or ${IMAGE_NAME},dynratelimiter)/dynratelimiter:$(or ${IMAGE_TAG},${COMMIT}) --file build/docker/dynratelimiter/Dockerfile .
	docker buildx build --push --platform linux/amd64,linux/i386,linux/arm64,linux/arm/v6,linux/arm/v7 -t $(or ${IMAGE_NAME},dynratelimiter)/dynratelimiter:latest --file build/docker/dynratelimiter/Dockerfile .

build-top:
	CGO_ENABLED=0 go build -o resource-top github.com/arivum/dynratelimiter/cmd/resource-top

build-stressme:
	CGO_ENABLED=0 go build -o stressme github.com/arivum/dynratelimiter/cmd/stressme

build-dynratelimiter-operator:
	CGO_ENABLED=0 go build -o dynratelimiter-operator github.com/arivum/dynratelimiter/cmd/dynratelimiter-operator

pack-dynratelimiter:
	docker build -t $(or ${IMAGE_NAME},dynratelimiter)/dynratelimiter:$(or ${IMAGE_TAG},${COMMIT}) --file build/docker/dynratelimiter/Dockerfile .
	docker build -t $(or ${IMAGE_NAME},dynratelimiter)/dynratelimiter:latest --file build/docker/dynratelimiter/Dockerfile .

pack-dynratelimiter-operator:
	docker build -t $(or ${IMAGE_NAME},dynratelimiter)/dynratelimiter-operator:$(or ${IMAGE_TAG},${COMMIT}) --file build/docker/dynratelimiter-operator/Dockerfile .
	docker build -t $(or ${IMAGE_NAME},dynratelimiter)/dynratelimiter-operator:latest --file build/docker/dynratelimiter-operator/Dockerfile .