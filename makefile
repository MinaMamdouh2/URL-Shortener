# Variables & dependencies
BASE_IMAGE_NAME := MinaMamdouh2/service
SERVICE_NAME    := url-shortener-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

run-local:
	./run-local.sh

run-local-air:
	./run-air.sh

# ==============================================================================
# Building containers
build-service-image:
	docker build \
		-f zarf/docker/dockerfile.service \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

docker-compose-up: build-service-image
	docker compose \
	-f "zarf\docker\docker-compose.service.yaml" up 