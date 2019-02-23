//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"testing"

	"github.com/nbio/st"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// TestRunPutCreditSpread01 - Run a put credit spread backtest.
//
func TestDoPutCreditSpread01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Setup a new backtesting
	bt := New(db)

	// Build screener object
	screen := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "short-strike-percent-away", Operator: "<", ValueNumber: 4.0},
			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.20},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// Run blank backtest
	err := bt.DoBacktestDays(models.Backtest{
		UserId:          1,
		StartingBalance: 10000.00,
		EndingBalance:   10000.00,
		StartDate:       helpers.ParseDateNoError("2018-01-01"),
		EndDate:         helpers.ParseDateNoError("2018-12-31"),
		Midpoint:        true,
		TradeSelect:     "highest-credit",
		Screen:          screen,
	})
	st.Expect(t, err, nil)

}

/* End File */
