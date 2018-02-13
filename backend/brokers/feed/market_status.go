//
// Date: 2/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import (
	"fmt"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

//
// Ticker - Get GetMarketStatus : 5 seconds
//
func (t *Base) DoGetMarketStatusTicker() {
	var err error

	for {

		// Load up market status.
		err = t.GetMarketStatus()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 10 second.
		time.Sleep(time.Second * 5)

	}
}

//
// Do get market status
//
func (t *Base) GetMarketStatus() error {
	// Make API call
	marketStatus, err := t.Api.GetMarketStatus()

	if err != nil {
		return err
	}

	// Save the market status in the fetch object
	t.muMarketStatus.Lock()
	t.MarketStatus = marketStatus
	t.muMarketStatus.Unlock()

	// Send up websocket.
	err = t.WriteDataChannel("market/status", marketStatus)

	if err != nil {
		return fmt.Errorf("GetMarketStatus() WriteDataChannel : ", err)
	}

	// Return Happy
	return nil
}

/* End File */
