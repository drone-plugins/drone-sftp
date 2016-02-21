.PHONY: all install test docker

IMAGE ?= plugins/drone-nuget

all: install test

install:
	npm install --quiet

test:
	@echo "Currently we don't provide test cases!"

docker:
	docker build --rm -t $(IMAGE) .
