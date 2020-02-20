#!/bin/bash

export PLUGIN_TEMPLATE=test/configmap.template.yaml
export PLUGIN_CONFIGMAP_FILE=test/sample-config-data
export PLUGIN_NAME=drone-kube-test
export PLUGIN_COMMIT=a5b81d0f
export PLUGIN_NAMESPACE=default

go build -o build/kubano
export $(cat .env | xargs) && ./build/kubano

# docker run --env-file=.env drone-kubano
