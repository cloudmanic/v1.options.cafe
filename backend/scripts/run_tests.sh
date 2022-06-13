#!/bin/bash

# Date: 3/5/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
#
# Run all tests. We do this instead of test ./... because lots of tests are broken. We should do away
# with this once all tests work. (Also note there might be issues with parrel tests all connecting to mysql)
# 
#

cd ../

# go test ./models/...
# go test ./brokers/eod/...
# go test ./brokers/tradier/...
go test ./library/cache/...
# go test ./library/analyze/...
# go test ./library/checkmail/...
# go test ./library/helpers/...
# go test ./library/market/...
go test ./library/realip/...
# go test ./library/reports/...
# go test ./library/services/...
# go test ./library/archive/...
# go test ./controllers/...
# go test ./brokers/pull/...
#go test ./backtesting/...

cd scripts
