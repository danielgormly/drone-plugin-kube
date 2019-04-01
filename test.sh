#!/bin/bash

export DRONE_COMMIT_SHA=1234567
export DRONE_BRANCH=test
export PLUGIN_CA=test
export PLUGIN_TOKEN=test
export PLUGIN_ENDPOINT=test
export PLUGIN_NAMESPACE=test
export PLUGIN_TEMPLATE=test/deployment.template.yaml
export PLUGIN_NAME=api

go build
./kubano
