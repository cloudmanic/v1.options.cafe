//
// Date: 2018-11-09
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-23
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package pull

import (
	"fmt"
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/archive"
	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/queue"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// DoGetAllOrders - Do get All orders. Here is where we archive historical orders.
//
func DoGetAllOrders(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {
	// Orders
	var orders []types.Order

	// Helpful log.
	services.Info("DoGetAllOrders() : Getting all orders for " + user.Email + ".")

	// Make API call
	orders, err := api.GetAllOrders()

	if err != nil {
		return fmt.Errorf("pull.GetAllOrders() : %s", err)
	}

	// Store the orders in our database
	err = archive.StoreOrders(db, orders, user.Id, broker.Id)

	if err != nil {
		return fmt.Errorf("pull.GetAllOrders() - StoreOrders() : %s", err)
	}

	// Consider this user boot strapped
	endBookstrapping(db, user.Id, user.Email)

	// Return Happy
	return nil
}

//
// Do get orders. Main thing we are doing here is populating the cache with the results
//
func DoGetOrders(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {

	// We do not call this until a user is boot strapped
	if user.Bootstrapped == "No" {
		services.Info("Skipping DoGetOrders() for user " + user.Email + " as the user is not bootstrapped yet.")
		return nil
	}

	// Make API call
	orders, err := api.GetOrders()

	if err != nil {
		return err
	}

	// Store result in cache.
	cache.Set("oc-orders-active-"+strconv.Itoa(int(user.Id))+"-"+strconv.Itoa(int(broker.Id)), orders)

	// Store symbols we use in orders
	for _, row := range orders {

		db.CreateActiveSymbol(user.Id, row.Symbol)

		for _, row2 := range row.Legs {
			db.CreateActiveSymbol(user.Id, row2.Symbol)
			db.CreateActiveSymbol(user.Id, row2.OptionSymbol)
		}

	}

	// Send message to websocket
	queue.Write("oc-websocket-write", `{"uri":"orders","user_id":`+strconv.Itoa(int(user.Id))+`,"body":`+helpers.JsonEncode(orders)+`}`)

	// Store the orders in our database
	err = archive.StoreOrders(db, orders, user.Id, broker.Id)

	if err != nil {
		return err
	}

	// Return Happy
	return nil
}

// ------------ Private Helper Functions -------------- //

//
// endBookstrapping - We just call this function to end the bookstraping process for a new account.
// We broke this out into a new function because we had an issue with the broker being set back to
// disabled. We need to update just one field not the entire user object.
//
func endBookstrapping(db models.Datastore, userID uint, email string) {
	// Update the DB record that we are now bootstrapped
	db.New().Model(&models.User{}).Where("id = ?", userID).Update("bootstrapped", "Yes")

	// Log action
	services.Info("DoGetAllOrders() : Setting user " + email + " to fully bootstrapped.")

	// Return Happy
	return
}

/* End File */
