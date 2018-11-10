//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import (
	"fmt"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/archive"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

//
// Ticker - Orders Archive : 24 hours
//
func (t *Base) DoOrdersTicker() {
	var err error
	var firstDone bool = false

	for {

		// Do we break out ?
		t.MuPolling.Lock()
		breakOut := t.Polling
		t.MuPolling.Unlock()

		if !breakOut {
			break
		}

		// Load up all orders
		_, err = t.GetAllOrders()

		if err != nil {
			services.Warning(err)
		}

		// Start the every 3 second ticker.
		// We do this because we want the archive
		// Action to run first to avoid race conditions.
		//
		// We also start the positions ticker here as we
		// want that to start after we archive all orders.
		if !firstDone {
			go t.DoPositionsTicker()
			firstDone = true
		}

		// Sleep for 24 hours
		time.Sleep(time.Hour * 24)
	}

	services.Info("Stopping DoOrdersTicker() : " + t.User.Email)
}

//
// Do get all orders. We return the orders instead of sending it up the websocket
//
func (t *Base) GetAllOrders() ([]types.Order, error) {

	var orders []types.Order

	// Make API call
	orders, err := t.Api.GetAllOrders()

	if err != nil {
		return orders, fmt.Errorf("Fetch.GetAllOrders() : ", err)
	}

	// Store the orders in our database
	err = archive.StoreOrders(t.DB, orders, t.User.Id, t.BrokerId)

	if err != nil {
		return orders, fmt.Errorf("Fetch.GetAllOrders() - StoreOrders() : ", err)
	}

	// Return Happy
	return orders, nil
}

/* End File */
