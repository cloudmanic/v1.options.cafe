//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/archive"
	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/notify/websocket_push"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cnf/structhash"
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
			go t.DoOrdersActiveTicker()
			go t.DoPositionsTicker()
			firstDone = true
		}

		// Sleep for 24 hours
		time.Sleep(time.Hour * 24)
	}

	services.Info("Stopping DoOrdersTicker() : " + t.User.Email)
}

//
// Ticker - Orders : 3 seconds
//
func (t *Base) DoOrdersActiveTicker() {
	var hash string = ""

	for {

		// Do we break out ?
		t.MuPolling.Lock()
		breakOut := t.Polling
		t.MuPolling.Unlock()

		if !breakOut {
			break
		}

		// Load up orders
		lastHash, err := t.GetOrders()

		if err != nil {
			services.Warning(err)
		}

		// If there has been any changes in our orders send a notice.
		if (len(hash) > 0) && (hash != lastHash) {
			websocket_push.Push(t.User.Id, "change-detected", `{ "type": "orders" }`)
			websocket_push.Push(t.User.Id, "change-detected", `{ "type": "trade-groups" }`)
		}

		// Store this hash for next time.
		hash = lastHash

		// Sleep for 3 second.
		time.Sleep(time.Second * 3)
	}

	services.Info("Stopping DoOrdersActiveTicker() : " + t.User.Email)
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

//
// Do get orders. Returns a hash of the orders
//
func (t *Base) GetOrders() (string, error) {

	orders := []types.Order{}

	// Make API call
	orders, err := t.Api.GetOrders()

	if err != nil {
		return "", err
	}

	// Store result in cache.
	cache.Set("oc-orders-active-"+strconv.Itoa(int(t.User.Id))+"-"+strconv.Itoa(int(t.BrokerId)), orders)

	// Store symbols we use in orders
	for _, row := range orders {

		t.DB.CreateActiveSymbol(t.User.Id, row.Symbol)

		for _, row2 := range row.Legs {
			t.DB.CreateActiveSymbol(t.User.Id, row2.Symbol)
			t.DB.CreateActiveSymbol(t.User.Id, row2.OptionSymbol)
		}

	}

	// Store the orders in our database
	err = archive.StoreOrders(t.DB, orders, t.User.Id, t.BrokerId)

	if err != nil {
		fmt.Errorf("Fetch.GetOrders() - StoreOrders() : ", err)
	}

	// Get a hash of the data structure.
	hash, err := structhash.Hash(orders, 1)

	if err != nil {
		return "", err
	}

	// Return Happy
	return hash, nil

}

/* End File */
