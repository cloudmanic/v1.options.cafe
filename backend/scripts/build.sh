#!/bin/bash

cd ../
# go build -o builds/options_cafe.darwin.amd64
env GOOS=linux GOARCH=amd64 go build -o builds/options_cafe.linux.amd64
cd scripts
