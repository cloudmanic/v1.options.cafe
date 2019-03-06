//
// Date: 2018-11-10
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-10
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package pull

import (
	"fmt"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/library/archive"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Do get history.
//
func DoGetHistory(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {

	history, err := api.GetAllHistory()

	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("pull.DoGetHistory(): %s", err)
	}

	// Store the history in our database
	err = archive.StoreBrokerEvents(db, history, user.Id, broker.Id)

	// Return Happy
	return nil
}

/* End File */
