//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import (
	"fmt"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Ticker - User Profile : 20 seconds
//
func (t *Base) DoUserProfileTicker() {

	var err error

	for {

		err = t.GetUserProfile()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 20 second.
		time.Sleep(time.Second * 20)

	}
}

//
// Fetch - Do get user profile
//
func (t *Base) GetUserProfile() error {

	// Make API call
	userProfile, err := t.Api.GetUserProfile()

	if err != nil {
		return err
	}

	// Get the broker defaults
	brokerConfig := t.Api.GetBrokerConfig()

	// Insert and/or Get Broker Account Record.
	for _, row := range userProfile.Accounts {

		// Create broker account object
		ba := &models.BrokerAccount{
			UserId:        t.User.Id,
			BrokerId:      t.BrokerId,
			AccountNumber: row.AccountNumber,
		}

		// Make DB query
		isNew, err := t.DB.FirstOrCreateBrokerAccount(ba)

		if err != nil {
			return fmt.Errorf("GetUserProfile() FirstOrCreateBrokerAccount : ", err)
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
			err := t.DB.UpdateBrokerAccount(ba)

			if err != nil {
				return fmt.Errorf("GetUserProfile() FirstOrCreateBrokerAccount : ", err)
			}

		}
	}

	// Save the orders in the fetch object
	t.muUserProfile.Lock()
	t.UserProfile = userProfile
	t.muUserProfile.Unlock()

	// Send up websocket.
	err = t.WriteDataChannel("user/profile", userProfile)

	if err != nil {
		return fmt.Errorf("GetUserProfile() WriteDataChannel : ", err)
	}

	// Return Happy
	return nil
}

/* End File */
