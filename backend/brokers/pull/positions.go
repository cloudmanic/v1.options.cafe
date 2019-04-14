//
// Date: 2018-11-10
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-23
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package pull

import (
	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/library/archive"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Fetch - Positions
//
func DoGetPositions(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {

	// We do not call this until a user is boot strapped
	if user.Bootstrapped == "No" {
		services.InfoMsg("Skipping DoGetPositions() for user " + user.Email + " as the user is not bootstrapped yet.")
		return nil
	}

	// Make API call
	positions, err := api.GetPositions()

	if err != nil {
		return err
	}

	// Loop through and add any positions to the active_symbols table.
	// Also see if we need to create a trade group
	for _, row := range positions {

		// Set active symbol
		db.CreateActiveSymbol(user.Id, row.Symbol)

		// Create a trade group if need be
		err = archive.PastCreateTradeGroupFromPosition(db, user.Id, broker.Id, row)

		if err != nil {
			services.Info(err)
			continue
		}
	}

	// Return Happy
	return nil
}

/* End File */
