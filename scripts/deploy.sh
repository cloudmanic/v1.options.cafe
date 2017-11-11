#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#
# Deploy the entire app (frontend and backend). We do some compiling locally and then deploy.

# Build the backend app within the docker container.
cd ../backend/docker

echo "Building app.options.cafe"
docker-compose run --rm app /work/scripts/build.sh

echo "Building cron.options.cafe"
docker-compose run --rm cron /work/scripts/build_cron.sh

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