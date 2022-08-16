//
// Date: 2022-07-14
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package user

import (
	"app.options.cafe/brokers/tradier"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
)

//
// ImportGainsForUser will make API calls to the broker and pull in profit / loss information and store in the database
//
func ImportGainsForUser(db *models.DB, brokerAccount models.BrokerAccount) {
	// Create new tradier instance
	tradier := &tradier.Api{ApiKey: "XXXXXXXXXXXXXXX"}

	// Run function
	gains, _ := tradier.GetGainsByAccountId(brokerAccount.AccountNumber)

	// Loop through the gains and store in our database.
	for _, row := range gains {
		// Figure out the symbol first
		var syb models.Symbol
		var err error

		// Option or stock
		if len(row.Symbol) > 10 {
			// Create options symbol
			syb, err = db.CreateNewOptionSymbol(row.Symbol)

			if err != nil {
				services.Fatal(err)
				continue
			}
		} else {
			// TODO(spicer): We should always have the symbol in our db. But we should check and deal with it.
			syb, err = db.CreateNewSymbol(row.Symbol, "", "Equity")

			if err != nil {
				services.Fatal(err)
				continue
			}
		}

		// Only add to the database if we do not have this record
		if db.New().Where("user_id = ? AND broker_account_id = ? AND symbol_id = ? AND open_date = ? AND close_date = ? AND cost = ? AND gain_loss = ? AND quantity = ?", brokerAccount.UserId, brokerAccount.Id, syb.Id, models.Date(row.OpenDate), models.Date(row.CloseDate), row.Cost, row.GainLoss, row.Quantity).First(&models.Gain{}).RecordNotFound() {
			// Model to insert
			i := models.Gain{
				UserId:          brokerAccount.UserId,
				BrokerAccountId: brokerAccount.Id,
				SymbolId:        syb.Id,
				Symbol:          syb,
				OpenDate:        models.Date(row.OpenDate),
				CloseDate:       models.Date(row.CloseDate),
				Cost:            row.Cost,
				GainLoss:        row.GainLoss,
				GainLossPercent: row.GainLossPercent,
				Proceeds:        row.Proceeds,
				Quantity:        row.Quantity,
				Term:            row.Term,
				Reviewed:        "No",
			}

			// Save to database
			db.New().Save(&i)
		}
	}

}
