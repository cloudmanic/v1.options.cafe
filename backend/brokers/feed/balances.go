//
// Date: 2/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import (
	"fmt"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/archive"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

//
// Ticker - Get GetBalances : 5 seconds
//
func (t *Base) DoGetBalancesTicker() {

	var err error
	count := 720

	for {

		// Do we break out ?
		t.MuPolling.Lock()
		breakOut := t.Polling
		t.MuPolling.Unlock()

		if !breakOut {
			break
		}

		// Load up balances. Every 3600 (1 hour) laps we log to the archive.
		if count > 720 {
			count = 0
			err = t.GetBalances(true)
		} else {
			count++
			err = t.GetBalances(false)
		}

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 5 second.
		time.Sleep(time.Second * 5)
	}

	services.Info("Stopping DoGetBalancesTicker() : " + t.User.Email)

}

//
// Do get Balances
//
func (t *Base) GetBalances(achive bool) error {
	balances, err := t.Api.GetBalances()

	if err != nil {
		return err
	}

	// Store balances in database.
	if achive {
		archive.StoreBalance(t.DB, balances, t.User.Id, t.BrokerId)
	}

	// Send up websocket.
	err = t.WriteDataChannel("balances", balances)

	if err != nil {
		return fmt.Errorf("GetBalances() WriteDataChannel : ", err)
	}

	// Return Happy
	return nil
}

/* End File */
