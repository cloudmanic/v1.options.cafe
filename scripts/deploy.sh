#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#
# Deploy the entire app (frontend and backend). We do some compiling locally and then deploy.

# Build backend
cd ../backend

echo "Building app.options.cafe"
env GOOS=linux GOARCH=amd64 go build -o builds/app.options.cafe

# Build frontend
cd ../frontend

echo "Building Frontend"
ng build -prod

cd ../scripts

# Deploy to backend with Ansible
cd ../ansible
ansible-playbook app.yml
cd ../scripts

# Login as myself and build and restart
ssh web2.cloudmanic.com "cd /sites/optionscafe/app.options.cafe/docker && docker-compose build && docker-compose down && docker-compose up -d && docker image prune -a -f"

## TODO: make an api call to papertail and output the current logs just to see if anything went wrong during deploy