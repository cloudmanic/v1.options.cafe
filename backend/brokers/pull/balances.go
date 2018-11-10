//
// Date: 2018-11-09
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-10
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package pull

import (
	"fmt"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Do get balances. Main thing we are doing here is populating the cache with the results
//
func DoGetBalances(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {

	fmt.Println(user.Email)

	// // Make API call
	// orders, err := api.GetOrders()

	// if err != nil {
	// 	return err
	// }

	// // Store result in cache.
	// cache.Set("oc-orders-active-"+strconv.Itoa(int(user.Id))+"-"+strconv.Itoa(int(broker.Id)), orders)

	// // Store symbols we use in orders
	// for _, row := range orders {

	// 	db.CreateActiveSymbol(user.Id, row.Symbol)

	// 	for _, row2 := range row.Legs {
	// 		db.CreateActiveSymbol(user.Id, row2.Symbol)
	// 		db.CreateActiveSymbol(user.Id, row2.OptionSymbol)
	// 	}

	// }

	// // Store the orders in our database
	// err = archive.StoreOrders(db, orders, user.Id, broker.Id)

	// if err != nil {
	// 	return err
	// }

	// Return Happy
	return nil
}

/* End File */
