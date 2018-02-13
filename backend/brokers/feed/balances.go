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
// Ticker - Get GetBalances : 5 seconds
//
func (t *Base) DoGetBalancesTicker() {
	var err error

	for {

		// Load up market status.
		err = t.GetBalances()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 5 second.
		time.Sleep(time.Second * 5)
	}
}

//
// Do get Balances
//
func (t *Base) GetBalances() error {
	balances, err := t.Api.GetBalances()

	if err != nil {
		return err
	}

	// Save the balances in the fetch object
	t.muBalances.Lock()
	t.Balances = balances
	t.muBalances.Unlock()

	// Send up websocket.
	err = t.WriteDataChannel("balances", balances)

	if err != nil {
		return fmt.Errorf("GetBalances() WriteDataChannel : ", err)
	}

	// Return Happy
	return nil
}

/* End File */
