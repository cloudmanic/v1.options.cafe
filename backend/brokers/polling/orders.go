//
// Date: 2018-11-09
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-10
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package polling

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Poll for orders by sending requests to our message queue.
// Every 3 seconds we send a request for more current order data.
//
func PollOrders(db models.Datastore) {

	// Start polling
	for {

		// Send action to all users
		SendActionToAllUsers(db, "get-orders")

		// Sleep for 3 seconds
		time.Sleep(time.Second * 3)
	}

}

/* End File */
