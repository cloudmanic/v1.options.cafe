//
// Date: 2/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package state

var (
	marketStatus interface{}
)

//
// Set market status.
//
func SetMarketStatus(status interface{}) {
	marketStatus = status
}

//
// Get market status.
//
func GetMarketStatus() interface{} {
	return marketStatus
}

/* End File */
