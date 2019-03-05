#!/bin/bash

# Date: 3/5/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
#
# Run all tests. We do this instead of test ./... because we use mysql and packages run in parallel.
# Downside to this is we need to remember to put each package into this by hand.
#

cd ../

#go test ./backtesting/...
go test ./brokers/eod/...
go test ./brokers/tradier/...
go test ./library/cache/...
go test ./library/analyze/...
go test ./library/checkmail/...
go test ./library/helpers/...
go test ./library/market/...
go test ./library/realip/...
go test ./library/reports/...
#go test ./library/services/...


cd scripts