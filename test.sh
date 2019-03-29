#!/bin/bash

export DRONE_BRANCH=test
go build
./kubano --server=go --token=lol --ca=ho --template=test/deployment.yaml
