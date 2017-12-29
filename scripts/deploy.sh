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

echo "Building cron.options.cafe"
cd cron
env GOOS=linux GOARCH=amd64 go build -o ../builds/cron.options.cafe

echo "Building cmd.options.cafe"
cd ../cmd
env GOOS=linux GOARCH=amd64 go build -o ../builds/cmd.options.cafe

# Build frontend
cd ../../frontend/docker

echo "Building Frontend"
docker-compose run --rm app ng build -prod

cd ../../scripts

# Deploy to backend with Ansible
cd ../ansible
ansible-playbook deploy.yml
cd ../scripts

# Login as myself and build and restart
ssh web2.cloudmanic.com "cd /sites/optionscafe/app.options.cafe/docker && docker-compose build && docker-compose down && docker-compose up -d"

## TODO: make an api call to papertail and output the current logs just to see if anything went wrong during deploy