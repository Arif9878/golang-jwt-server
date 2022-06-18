HUB ?= docker.io/istio/jwt-server
TAG ?= 0.6

run:
	@go run cmd/main.go

build: main.go Dockerfile
	docker build . -t $(HUB):$(TAG)

push: build
	docker push $(HUB):$(TAG)
