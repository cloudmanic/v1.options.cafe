//
// Date: 2018-11-09
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-11
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package pull

import (
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/library/archive"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/queue"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Do get balances. Main thing we are doing here is populating the cache with the results
//
func DoGetBalances(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {

	balances, err := api.GetBalances()

	if err != nil {
		return err
	}

	// Store balances in database.
	go archive.StoreBalance(db, balances, user.Id, broker.Id)

	// Send message to websocket
	queue.Write("oc-websocket-write", `{"uri":"balances","user_id":`+strconv.Itoa(int(user.Id))+`,"body":`+helpers.JsonEncode(balances)+`}`)

	// Return Happy
	return nil
}

/* End File */
