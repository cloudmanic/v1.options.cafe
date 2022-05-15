//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"testing"

	"app.options.cafe/brokers/eod"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
	screenerCache "app.options.cafe/screener/cache"
	"github.com/davecgh/go-spew/spew"
	"github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
)

//
// TestLongCallButterflySpreadResults01 - Test finding results
//
func TestLongCallButterflySpreadResults01(t *testing.T) {
	// Only do this for non-short
	if testing.Short() {
		t.Skipf("Skipping TestLongCallButterflySpreadResults01 test since it requires a env tokens and --short was requested")
	}

	// Load .env file (MUST CAll GO TEST FROM THE ROOT)
	env.ReadEnv("../.env")

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Vars needed for backtest
	today := helpers.ParseDateNoError("2022-03-10")

	// Build screener object
	screen := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Name:     "SPY Long Call Butterfly",
		Strategy: "long-call-butterfly-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "left-strike-percent-away", Operator: ">", ValueNumber: 3},
			{UserId: 1, Key: "right-strike-percent-away", Operator: ">", ValueNumber: 3},
			{UserId: 1, Key: "open-debit", Operator: "<", ValueNumber: 3.00},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 31},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
			// {UserId: 1, Key: "allow-more-than-one-expire", Operator: "=", ValueString: "no"},
			// {UserId: 1, Key: "allow-more-than-one-strike", Operator: "=", ValueString: "no"},
		},
	}

	// Set backtest
	backtest := models.Backtest{
		UserId:          uint(screen.UserId),
		StartingBalance: 5000.00,
		EndingBalance:   5000.00,
		PositionSize:    "10-percent", // one-at-time, *-percent
		StartDate:       models.Date{helpers.ParseDateNoError("2022-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2022-12-31")},
		Midpoint:        true,
		//TradeSelect:     "shortest-percent-away", // least-days-to-expire, highest-midpoint, highest-ask, highest-percent-away, shortest-percent-away
		Benchmark: "SPY",
		Screen:    screen,
	}

	// Create broker object
	o := eod.Api{
		DB:  db,
		Day: today,
	}

	// Get all options for this symbol and day.
	options, underlyingLast, err := o.GetOptionsBySymbol(backtest.Screen.Symbol)

	if err != nil {
		services.Info(err)
	}

	// Setup cache
	cache := screenerCache.New(db)

	// Call the function we are testing.
	bt := New(db, 1, backtest.Benchmark)
	screenerResults, err := bt.LongCallButterflySpreadResults(today, &backtest, underlyingLast, options, cache)
	st.Expect(t, err, nil)

	spew.Dump(screenerResults)
}
