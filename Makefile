.PHONY: install test docker

IMAGE ?= plugins/drone-nuget

install:
	npm install --quiet

test:
	@echo "Currently we don't provide test cases!"

docker:
	docker build --rm -t $(IMAGE) .
