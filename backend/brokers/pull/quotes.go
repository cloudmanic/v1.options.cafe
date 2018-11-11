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
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/queue"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Do get quotes.
//
func DoGetQuotes(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {

	var symbols []string

	// Build active symbols array.
	results, err := db.GetActiveSymbolsByUser(user.Id)

	for _, row := range results {
		symbols = append(symbols, row.Symbol)
	}

	// Api call to get quote
	detailedQuotes, err := api.GetQuotes(symbols)

	if err != nil {
		return err
	}

	// Loop through the quotes sending them up the websocket channel
	for _, row := range detailedQuotes {

		// Send message to websocket
		queue.Write("oc-websocket-write", `{"uri":"quote","user_id":`+strconv.Itoa(int(user.Id))+`,"body":`+helpers.JsonEncode(row)+`}`)

	}

	// Return happy
	return nil
}

/* End File */
