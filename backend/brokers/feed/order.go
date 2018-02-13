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
func (t *Base) DoOrdersArchive() {
	var err error

	for {
		// Load up all orders
		_, err = t.GetAllOrders()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 24 hours
		time.Sleep(time.Hour * 24)
	}
}

//
// Ticker - Orders : 3 seconds
//
func (t *Base) DoOrdersTicker() {
	var err error

	for {
		// Load up orders
		err = t.GetOrders()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 3 second.
		time.Sleep(time.Second * 3)
	}
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
// Do get orders
//
func (t *Base) GetOrders() error {

	orders := []types.Order{}

	// Make API call
	orders, err := t.Api.GetOrders()

	if err != nil {
		return err
	}

	// Save the orders in the fetch object
	t.muOrders.Lock()
	t.Orders = orders
	t.muOrders.Unlock()

	// Store the orders in our database
	err = archive.StoreOrders(t.DB, orders, t.User.Id, t.BrokerId)

	if err != nil {
		fmt.Errorf("Fetch.GetOrders() - StoreOrders() : ", err)
	}

	// Send up websocket.
	err = t.WriteDataChannel("orders", orders)

	if err != nil {
		return fmt.Errorf("Fetch.GetOrders() : ", err)
	}

	// Return Happy
	return nil

}

/* End File */
