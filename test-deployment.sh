#!/bin/bash

PLUGIN_TEMPLATE=test/deployment.template.yaml
PLUGIN_NAME=drone-kube-test
PLUGIN_COMMIT=a5b81d0f

go build -o build/kubano
export $(cat .env | xargs) && ./build/kubano

# docker run --env-file=.env drone-kubano
