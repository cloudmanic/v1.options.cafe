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

	"app.options.cafe/brokers"
	"app.options.cafe/models"
)

//
// DoGetUserProfile - Do get user profile data.
//
func DoGetUserProfile(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {

	// Make API call
	userProfile, err := api.GetUserProfile()

	if err != nil {
		return err
	}

	// Get the broker defaults
	brokerConfig := api.GetBrokerConfig()

	// Insert and/or Get Broker Account Record.
	for _, row := range userProfile.Accounts {

		// Create broker account object
		ba := &models.BrokerAccount{
			UserId:        user.Id,
			BrokerId:      broker.Id,
			AccountNumber: row.AccountNumber,
		}

		// Make DB query
		isNew, err := db.FirstOrCreateBrokerAccount(ba)

		if err != nil {
			return fmt.Errorf("DoGetUserProfile() FirstOrCreateBrokerAccount : %s", err)
		}

		// If this is a new entry we should add default commissions.
		if isNew {

			ba.Name = ba.AccountNumber
			ba.StockCommission = brokerConfig.DefaultStockCommission
			ba.StockMin = brokerConfig.DefaultStockMin
			ba.OptionCommission = brokerConfig.DefaultOptionCommission
			ba.OptionSingleMin = brokerConfig.DefaultOptionSingleMin
			ba.OptionMultiLegMin = brokerConfig.DefaultOptionMultiLegMin
			ba.OptionBase = brokerConfig.DefaultOptionBase
			err := db.UpdateBrokerAccount(ba)

			if err != nil {
				return fmt.Errorf("DoGetUserProfile() FirstOrCreateBrokerAccount : %s", err)
			}

		}
	}

	// Return Happy
	return nil
}

/* End File */
