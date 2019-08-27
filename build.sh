#!/bin/bash

rm -rf build/kubano
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/kubano

docker build -t danielgormly/drone-plugin-kube:0.0.2 build
docker push danielgormly/drone-plugin-kube:0.0.2
