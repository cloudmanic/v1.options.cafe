//
// Date: 2018-09-14
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-09-14
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"errors"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Pass in all historical balances. One balance snapshot per day.
//
func StoreBalance(db models.Datastore, balances []types.Balance, userId uint, brokerId uint) error {

	// Loop through the balances
	for _, row := range balances {

		// Get broker account id
		brokerAccount := models.BrokerAccount{}
		db.New().Where("broker_id = ? AND account_number = ?", brokerId, row.AccountNumber).First(&brokerAccount)

		if brokerAccount.Id <= 0 {
			services.Info(errors.New("Broker account not found - " + row.AccountNumber))
			continue
		}

		// See if we already have this in our system
		bh := models.BalanceHistory{}
		db.New().Where("account_number = ? AND broker_account_id = ? AND date = ?", row.AccountNumber, brokerAccount.Id, models.Date{time.Now()}).First(&bh)

		if bh.Id > 0 {

			bh.AccountNumber = row.AccountNumber
			bh.AccountValue = row.AccountValue
			bh.TotalCash = row.TotalCash
			bh.OptionBuyingPower = row.OptionBuyingPower
			bh.StockBuyingPower = row.StockBuyingPower
			db.New().Save(&bh)

		} else {

			// Build object
			history := models.BalanceHistory{
				UserId:            userId,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				BrokerAccountId:   brokerAccount.Id,
				Date:              models.Date{time.Now()},
				AccountNumber:     row.AccountNumber,
				AccountValue:      row.AccountValue,
				TotalCash:         row.TotalCash,
				OptionBuyingPower: row.OptionBuyingPower,
				StockBuyingPower:  row.StockBuyingPower,
			}

			// Save record
			err := db.CreateNewRecord(&history, models.InsertParam{})

			if err != nil {
				services.Fatal(err)
			}

		}
	}

	// Return happy
	return nil
}

/* End File */
