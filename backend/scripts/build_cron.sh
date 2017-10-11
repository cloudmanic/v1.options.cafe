#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#
# This should be run from the docker container.

cd /work/cron

# Make sure we have everything we need with packages
go get

# Build for this docker image.
go build -o /work/builds/cron.options.cafe

cd /work/scripts