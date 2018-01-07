#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#
# Deploy backups for the app.

# Deploy to backups with Ansible
cd ../ansible
ansible-playbook backups.yml
cd ../scripts

# Login as myself and build and restart
ssh web2.cloudmanic.com "cd /sites/optionscafe/backup.options.cafe && docker-compose build && docker-compose down && docker-compose up -d"

## TODO: make an api call to papertail and output the current logs just to see if anything went wrong during deploy