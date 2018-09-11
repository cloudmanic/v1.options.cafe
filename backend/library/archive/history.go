//
// Date: 2018-09-11
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-09-11
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"strings"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Pass in all historical events and store them in our database.
//
func StoreHistory(db models.Datastore, history []types.History, userId uint, brokerId uint) error {

	// Loop through the history and process
	for _, row := range history {

		// Get broker account id
		brokerAccount, err := db.GetBrokerAccountByBrokerAccountNumber(brokerId, row.BrokerId)

		if err != nil {
			services.Fatal(err)
			continue
		}

		// See if we already have this in our system
		var ev models.BrokerEvent
		db.New().Where("broker_id = ? AND broker_account_id = ?", row.Id, brokerAccount.Id).First(&ev)

		if ev.Id > 0 {
			continue
		}

		// Set the Type
		eventType := strings.Title(row.Type)

		found, _ := helpers.InArray(eventType, []string{"Ach", "Trade", "Option", "Interest", "Journal", "Dividend", "Adjustment", "Other"})

		if !found {
			eventType = "Other"
		}

		// Set the Trade Type
		tradeType := strings.Title(row.TradeType)

		found2, _ := helpers.InArray(tradeType, []string{"Equity", "Option", "ETF", "Preferred Stock", "Other"})

		if !found2 {
			tradeType = "Other"
		}

		// Build object
		event := models.BrokerEvent{
			UserId:          userId,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			BrokerAccountId: brokerAccount.Id,
			BrokerId:        row.Id,
			Type:            eventType,
			Date:            helpers.ParseDateNoError(row.Date),
			Amount:          row.Amount,
			Symbol:          row.Symbol,
			Commission:      row.Commission,
			Description:     row.Description,
			Price:           row.Price,
			Quantity:        row.Quantity,
			TradeType:       tradeType,
		}

		// Save record
		err = db.CreateNewRecord(&event, models.InsertParam{})

		if err != nil {
			services.Fatal(err)
		}
	}

	// Return happy
	return nil
}

/* End File */
