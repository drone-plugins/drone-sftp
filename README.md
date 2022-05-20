# drone-sftp

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-sftp/status.svg)](http://cloud.drone.io/drone-plugins/drone-sftp)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-sftp?status.svg)](http://godoc.org/github.com/drone-plugins/drone-sftp)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-sftp)](https://goreportcard.com/report/github.com/drone-plugins/drone-sftp)

Drone plugin to publish files and artifacts via SFTP. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

go build -v -a -tags netgo -o release/linux/amd64/drone-sftp
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/sftp .
```

## Usage

```console
docker run --rm \
  -e PLUGIN_HOST=sftp.company.com \
  -e PLUGIN_PORT=2222 \
  -e PLUGIN_USERNAME=user \
  -e PLUGIN_PASSWORD=pa$$word \
  -e PLUGIN_FILES=*.nupkg
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/sftp
```
