.PHONY: build-image

build-image:
	@echo "build image..."
	docker build -t mxtransporter -f Dockerfile .

build-image-for-local:
	@echo "build image..."
	docker build -t mxtransporter -f Dockerfile.local .
