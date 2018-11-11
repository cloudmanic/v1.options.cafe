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
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// See if we need to refresh access tokens
//
func DoAccessTokenRefresh(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {

	fmt.Println(user.Email)

	err := api.DoRefreshAccessTokenIfNeeded(user)

	if err != nil {
		return err
	}

	// Return Happy
	return nil

}

/* End File */
