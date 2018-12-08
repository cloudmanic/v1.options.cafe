#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#
# Deploy the entire app (frontend and backend). We do some compiling locally and then deploy.

# Build backend
cd ../backend

echo "Building app.options.cafe"
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/app.options.cafe
upx builds/app.options.cafe 

# Build frontend
cd ../frontend

echo "Building Frontend"
ng build -prod

cd ../scripts

# Deploy to backend with Ansible
cd ../ansible
ansible-playbook app.yml
cd ../scripts

## TODO: make an api call to papertail and output the current logs just to see if anything went wrong during deploy