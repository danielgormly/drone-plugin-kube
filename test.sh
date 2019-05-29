#!/bin/bash

go build -o build/kubano
export $(cat .env | xargs) && ./build/kubano

# docker run --env-file=.env drone-kubano
