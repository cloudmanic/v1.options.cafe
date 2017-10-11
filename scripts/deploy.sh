#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#
# Deploy the entire app (frontend and backend). We do some compiling locally and then deploy.

# Build the backend app within the docker container.
cd ../backend/docker
echo "Building app.options.cafe"
docker-compose run --rm app scripts/build.sh
cd ../../scripts

# Deploy to backend with Ansible
cd ../ansible
ansible-playbook deploy.yml
cd ../scripts
