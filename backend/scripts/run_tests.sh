#!/bin/bash

# Date: 3/5/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
#
# Deploy the entire app (frontend and backend). We do some compiling locally and then deploy.

cd ../

go test ./users/...
go test ./models/...
go test ./library/...
go test ./brokers/...
go test ./websocket/...
go test ./controllers/...

cd scripts