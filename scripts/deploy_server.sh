#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#

# Deploy to backend with Ansible
cd ../ansible
ansible-playbook server.yml
cd ../scripts

# Login as myself and build and restart
ssh web2.cloudmanic.com "cd /sites/optionscafe/app.options.cafe/docker && docker-compose build && docker-compose down && docker-compose up -d"

# Deploy app
./deploy.sh

## TODO: make an api call to papertail and output the current logs just to see if anything went wrong during deploy