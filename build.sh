#!/usr/bin/env bash

set -eou pipefail

rm -rf build/

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/kubano

docker build -t polygonio/sandbox:drone-plugin-kube build
docker push polygonio/sandbox:drone-plugin-kube
