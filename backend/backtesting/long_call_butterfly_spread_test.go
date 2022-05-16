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
		TradeSelect:     "lowest-ask",
		Benchmark:       "SPY",
		Screen:          screen,
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
	st.Expect(t, len(screenerResults), 4)
	st.Expect(t, screenerResults[0].Debit, 2.84)
	st.Expect(t, screenerResults[0].MidPoint, 0.27)
	st.Expect(t, len(screenerResults[0].Legs), 3)
}

//
// TestLongCallButterflySpreadPlaceTrades01 - Selects what trades to make - lowest-ask
//
func TestLongCallButterflySpreadPlaceTrades01(t *testing.T) {
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
		StartingBalance: 50000.00,
		EndingBalance:   50000.00,
		PositionSize:    "10-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2022-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2022-12-31")},
		Midpoint:        false,
		TradeSelect:     "lowest-ask",
		Benchmark:       "SPY",
		Screen:          screen,
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
	st.Expect(t, len(screenerResults), 4)
	st.Expect(t, screenerResults[0].Debit, 2.84)
	st.Expect(t, screenerResults[0].MidPoint, 0.27)
	st.Expect(t, len(screenerResults[0].Legs), 3)

	// Test the function we are here to test.
	bt.LongCallButterflySpreadPlaceTrades(today, &backtest, screenerResults, options)
	st.Expect(t, len(backtest.TradeGroups), 1)
	st.Expect(t, backtest.TradeGroups[0].Lots, 17)
	st.Expect(t, backtest.TradeGroups[0].Positions[0].Qty, 17)
	st.Expect(t, backtest.TradeGroups[0].Positions[1].Qty, -34)
	st.Expect(t, backtest.TradeGroups[0].Positions[2].Qty, 17)
	st.Expect(t, backtest.TradeGroups[0].Positions[0].Symbol.ShortName, "SPY220330C00412000")
	st.Expect(t, backtest.TradeGroups[0].Positions[1].Symbol.ShortName, "SPY220330C00425000")
	st.Expect(t, backtest.TradeGroups[0].Positions[2].Symbol.ShortName, "SPY220330C00438000")
	st.Expect(t, backtest.TradeGroups[0].SpreadText, "SPY 03/30/2022 $412.00 / $425.00 / $438.00")
}

//
// TestLongCallButterflySpreadPlaceTrades02 - Selects what trades to make - lowest-ask - Testing One trade at a time.
//
func TestLongCallButterflySpreadPlaceTrades02(t *testing.T) {
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
	tomorrow := helpers.ParseDateNoError("2022-03-11")

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
		StartingBalance: 50000.00,
		EndingBalance:   50000.00,
		PositionSize:    "one-at-time",
		StartDate:       models.Date{helpers.ParseDateNoError("2022-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2022-12-31")},
		Midpoint:        false,
		TradeSelect:     "lowest-ask",
		Benchmark:       "SPY",
		Screen:          screen,
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
	st.Expect(t, len(screenerResults), 4)
	st.Expect(t, screenerResults[0].Debit, 2.84)
	st.Expect(t, screenerResults[0].MidPoint, 0.27)
	st.Expect(t, len(screenerResults[0].Legs), 3)

	// Call the function we are testing.
	bt2 := New(db, 1, backtest.Benchmark)
	screenerResults2, err := bt2.LongCallButterflySpreadResults(tomorrow, &backtest, underlyingLast, options, cache)
	st.Expect(t, err, nil)
	st.Expect(t, len(screenerResults), 4)
	st.Expect(t, screenerResults[0].Debit, 2.84)
	st.Expect(t, screenerResults[0].MidPoint, 0.27)
	st.Expect(t, len(screenerResults[0].Legs), 3)

	// Test the function we are here to test.
	bt.LongCallButterflySpreadPlaceTrades(today, &backtest, screenerResults, options)
	st.Expect(t, len(backtest.TradeGroups), 1)

	// Test the function we are here to test. (nothing should change because we already have a trade on.)
	bt2.LongCallButterflySpreadPlaceTrades(tomorrow, &backtest, screenerResults2, options)
	st.Expect(t, len(backtest.TradeGroups), 1)
	st.Expect(t, backtest.TradeGroups[0].Lots, 1)
	st.Expect(t, backtest.TradeGroups[0].Positions[0].Qty, 1)
	st.Expect(t, backtest.TradeGroups[0].Positions[1].Qty, -2)
	st.Expect(t, backtest.TradeGroups[0].Positions[2].Qty, 1)
	st.Expect(t, backtest.TradeGroups[0].Positions[0].Symbol.ShortName, "SPY220330C00412000")
	st.Expect(t, backtest.TradeGroups[0].Positions[1].Symbol.ShortName, "SPY220330C00425000")
	st.Expect(t, backtest.TradeGroups[0].Positions[2].Symbol.ShortName, "SPY220330C00438000")
	st.Expect(t, backtest.TradeGroups[0].SpreadText, "SPY 03/30/2022 $412.00 / $425.00 / $438.00")
}

//
// TestCloseOnDebit01 - Close trade if we hit our debit goals
//
func TestCloseOnDebit01(t *testing.T) {
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
	today := helpers.ParseDateNoError("2022-04-11")
	testDay := helpers.ParseDateNoError("2022-04-19")

	// Build screener object
	screen := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Name:     "SPY Long Call Butterfly",
		Strategy: "long-call-butterfly-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "left-strike-percent-away", Operator: ">", ValueNumber: 2},
			{UserId: 1, Key: "right-strike-percent-away", Operator: ">", ValueNumber: 2},
			{UserId: 1, Key: "open-debit", Operator: "<", ValueNumber: 3.00},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 10},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
			{UserId: 1, Key: "take-profit-percent", Operator: ">", ValueNumber: 50},
		},
	}

	// Set backtest
	backtest := models.Backtest{
		UserId:          uint(screen.UserId),
		StartingBalance: 5000.00,
		EndingBalance:   5000.00,
		PositionSize:    "10-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2022-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2022-12-31")},
		Midpoint:        false,
		TradeSelect:     "lowest-ask",
		Benchmark:       "SPY",
		Screen:          screen,
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

	// Get the benchmark last
	backtest.BenchmarkStart = bt.getBenchmarkByDate(today)

	screenerResults, err := bt.LongCallButterflySpreadResults(today, &backtest, underlyingLast, options, cache)
	st.Expect(t, err, nil)

	// Test the function we are here to test.
	bt.LongCallButterflySpreadPlaceTrades(today, &backtest, screenerResults, options)
	st.Expect(t, len(backtest.TradeGroups), 1)

	// Get all options for this symbol
	o2 := eod.Api{
		DB:  db,
		Day: testDay,
	}

	options2, _, err := o2.GetOptionsBySymbol(backtest.Screen.Symbol)

	if err != nil {
		services.Info(err)
	}

	// Close out the trade
	bt.CloseOnDebit(testDay, &backtest, options2)
	st.Expect(t, backtest.EndingBalance, 5214.00)
	st.Expect(t, backtest.TradeGroups[0].Status, "Closed")
	st.Expect(t, backtest.TradeGroups[0].OpenPrice, 287.00)
	st.Expect(t, backtest.TradeGroups[0].ClosePrice, 501.00)
	st.Expect(t, backtest.TradeGroups[0].Balance, 5214.00)
	st.Expect(t, backtest.TradeGroups[0].ReturnPercent, 74.56)
	st.Expect(t, backtest.TradeGroups[0].BenchmarkLast, 445.04)
	st.Expect(t, backtest.TradeGroups[0].BenchmarkBalance, 5056.32)
	st.Expect(t, backtest.TradeGroups[0].BenchmarkReturn, 1.16)

	//spew.Dump(backtest)
}
