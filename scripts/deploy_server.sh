#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#

# Deploy to backend with Ansible
cd ../ansible
ansible-playbook server.yml
cd ../scripts

# Deploy app
./deploy.sh

## TODO: make an api call to papertail and output the current logs just to see if anything went wrong during deploy