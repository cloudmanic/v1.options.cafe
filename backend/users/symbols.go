//
// Date: 11/07/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import (
	"encoding/json"

	"app.options.cafe/backend/brokers/types"
	"app.options.cafe/backend/controllers"
	"app.options.cafe/backend/library/services"
	"github.com/tidwall/gjson"
)

//
// Pass in a query query and return a list of companies that match that query.
//
func (t *Base) SearchBySymbolOrCompanyName(user *UserFeed, request controllers.ReceivedStruct) {

	// Get query string
	query := gjson.Get(request.Body, "body.query").String()

	// Search for this query
	syms, err := user.BrokerFeed[1].Api.SearchBySymbolOrCompanyName(query)

	if err != nil {
		services.Error(err, "SearchBySymbolOrCompanyName() API call.")
		return
	}

	// We really only want to return the top 10 results.
	count := 0
	symbols := []types.Symbol{}
	for _, row := range syms {

		symbols = append(symbols, row)

		count++

		if count >= 5 {
			break
		}
	}

	// Convert to a json string.
	dataJson, err := json.Marshal(symbols)

	if err != nil {
		services.Error(err, "SearchBySymbolOrCompanyName() json.Marshal (#1)")
		return
	}

	// Build JSON we send
	jsonSend, err := t.WsSendJsonBuild("symbols/search", dataJson)

	if err != nil {
		services.Error(err, "SearchBySymbolOrCompanyName() WsSendJsonBuild (#2)")
		return
	}

	// Send back just to this particular connection.
	request.Connection.WriteChan <- jsonSend
}

/* End File */
