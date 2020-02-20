#!/bin/bash

export PLUGIN_TEMPLATE=test/deployment.template.yaml
export PLUGIN_NAME=drone-kube-test
export PLUGIN_NAMESPACE=default

go build -o build/kubano
export $(cat .env | xargs) && ./build/kubano

# docker run --env-file=.env drone-kubano
