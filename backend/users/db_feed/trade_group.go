//
// Date: 2/21/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package db_feed

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cnf/structhash"
)

//
// Start a feed to send to clients for
// Trade groups. Query the db on a timer
// and send the data to the front end.
//
func (t *Feed) DoTradeGroupsOpen() {

	var lastHash string = ""

	for {

		// Run the query - We limit to 25 as it is a pretty structured dataset.
		results, _, err := t.DB.GetTradeGroups(models.QueryParam{
			UserId:   t.User.Id,
			Order:    "id",
			Sort:     "asc",
			PreLoads: []string{"Positions"},
			Wheres: []models.KeyValue{
				{Key: "status", Value: "Open"},
			},
		})

		if err != nil {
			services.BetterError(err)
			continue
		}

		// Review the results for a particular hd5 hash.
		hash, err := structhash.Hash(results, 1)

		if err != nil {
			services.BetterError(err)
			continue
		}

		// See if we need to send a notice of a change
		if (len(hash) > 0) && (len(lastHash) > 0) && (hash != lastHash) {

			// Log the change.
			services.Info("A TradeGroup change has been detected for : " + t.User.Email)

			// Send the change up the websocket.
			t.WsSend("change-detected", `{ "type": "trade-groups-open" }`)
		}

		// Store hash for the next lap
		lastHash = hash

		// Sleep for 1 second.
		time.Sleep(time.Second * 1)
	}

}

/* End File */
