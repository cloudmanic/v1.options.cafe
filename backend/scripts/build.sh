#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#
# This should be run from the docker container.

cd /work

# Make sure we have everything we need with packages
go get

# Build for this docker image.
go build -o /work/builds/app.options.cafe

# Build for OSX
# go build -o builds/options_cafe.darwin.amd64

# Build for Linux server
#env GOOS=linux GOARCH=amd64 go build -o builds/options_cafe.linux.amd64

cd /work/scripts
