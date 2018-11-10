//
// Date: 2018-11-10
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
// Poll for quotes by sending requests to our message queue.
// Every 1 second we send a request for more current quote data.
//
func PollQuotes(db models.Datastore) {

	// Start polling
	for {

		// Send action to all users
		SendActionToAllUsers(db, "get-quotes")

		// Sleep for 1 seconds
		time.Sleep(time.Second * 1)
	}

}

/* End File */
