#!/bin/bash

# Date: 3/4/2018
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
#
# Start a docker image just for unit testing OC databases.
#

docker run --name options_cafe_testing -e MYSQL_ROOT_PASSWORD=foobar --tmpfs /var/lib/mysql -p 127.0.0.1:9906:3306 -d mariadb:10.2
