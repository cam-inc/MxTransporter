.PHONY: build-image

build-image:
	@echo "build image..."
	docker build -t mxtransporter .