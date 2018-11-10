//
// Date: 2018-11-09
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-09
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package polling

import (
	"os"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	nsq "github.com/nsqio/go-nsq"
)

//
// Poll for orders by sending requests to our message queue.
// Every 3 seconds we send a request for more current order data.
//
func PollOrders(db models.Datastore) {

	// NSQ config
	config := nsq.NewConfig()

	// Build producer object
	w, err := nsq.NewProducer(os.Getenv("NSQ_HOST"), config)

	if err != nil {
		services.FatalMsg(err, "PollOrders (NewProducer): NSQ Could not connect.")
	}

	// Close on function exit
	defer w.Stop()

	// Start polling
	for {

		// Send action to all users
		SendActionToAllUsers(db, "get-orders")

		// Sleep for 3 seconds
		time.Sleep(time.Second * 3)
	}

}

/* End File */
