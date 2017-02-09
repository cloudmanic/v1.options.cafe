#!/bin/bash

cd ../

# Build for OSX
# go build -o builds/options_cafe.darwin.amd64
# upx builds/options_cafe.darwin.amd64

# Build for Linux server
env GOOS=linux GOARCH=amd64 go build -o builds/options_cafe.linux.amd64
upx builds/options_cafe.linux.amd64

cd scripts
