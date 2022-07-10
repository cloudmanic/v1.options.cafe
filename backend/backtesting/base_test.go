//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"testing"

	"app.options.cafe/brokers/eod"
	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
	"github.com/nbio/st"

	"github.com/jpfuentes2/go-env"
)

//
// TestDoBacktestDays01 - Run a backtest looping through each day
//
func TestDoBacktestDays01(t *testing.T) {
	// Only do this for non-short
	if testing.Short() {
		t.Skipf("Skipping TestDoBacktestDays01 test since it requires a env tokens and --short was requested")
	}

	// Load .env file (MUST CAll GO TEST FROM THE ROOT)
	env.ReadEnv("../.env")

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Setup a new backtesting
	bt := New(db, 1, "SPY")

	// Build screener object
	screen := models.Screener{
		Symbol:   "SPY",
		Strategy: "empty",
	}

	// Run blank backtest
	err := bt.DoBacktestDays(&models.Backtest{
		StartDate: models.Date{helpers.ParseDateNoError("2020-01-01")},
		EndDate:   models.Date{helpers.ParseDateNoError("2020-01-10")},
		Screen:    screen,
	})
	st.Expect(t, err, nil)
}

//
// TestGetOptionsByExpirationType01
//
func TestGetOptionsByExpirationType01(t *testing.T) {
	// Only do this for non-short
	if testing.Short() {
		t.Skipf("Skipping TestGetOptionsByExpirationType01 test since it requires a env tokens and --short was requested")
	}

	// Load .env file (MUST CAll GO TEST FROM THE ROOT)
	env.ReadEnv("../.env")

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Setup a new backtesting
	bt := New(db, 1, "SPY")

	// Setup EOD Api
	o := eod.Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-01-03"),
	}

	// Get a list of options
	options, underlyingLast, err := o.GetOptionsBySymbol("spy")

	// Double check results.
	st.Expect(t, err, nil)
	st.Expect(t, len(options), 4512)
	st.Expect(t, 270.47, underlyingLast)

	// Get the options and just pull out the PUT options for this expire date.
	putOptions := bt.GetOptionsByExpirationType(types.Date{o.Day}, "Put", options)

	// Test results - Puts
	st.Expect(t, len(putOptions), 42)
	st.Expect(t, putOptions[10].OptionType, "Put")
	st.Expect(t, putOptions[11].OptionType, "Put")
	st.Expect(t, putOptions[21].OptionType, "Put")
	st.Expect(t, putOptions[30].OptionType, "Put")

	// Get the options and just pull out the Call options for this expire date.
	callOptions := bt.GetOptionsByExpirationType(types.Date{o.Day}, "Call", options)

	// Test results - Puts
	st.Expect(t, len(callOptions), 42)
	st.Expect(t, callOptions[10].OptionType, "Call")
	st.Expect(t, callOptions[11].OptionType, "Call")
	st.Expect(t, callOptions[21].OptionType, "Call")
	st.Expect(t, callOptions[30].OptionType, "Call")

	// This should error - rather return nothing - Note call vs. Call
	errorOptions := bt.GetOptionsByExpirationType(types.Date{o.Day}, "call", options)

	// Test results - Error
	st.Expect(t, len(errorOptions), 0)
}

//
// TestGetExpirationDatesFromOptions
//
func TestGetExpirationDatesFromOptions01(t *testing.T) {
	// Only do this for non-short
	if testing.Short() {
		t.Skipf("Skipping TestGetExpirationDatesFromOptions01 test since it requires a env tokens and --short was requested")
	}

	// Load .env file (MUST CAll GO TEST FROM THE ROOT)
	env.ReadEnv("../.env")

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Setup a new backtesting
	bt := New(db, 1, "SPY")

	// Setup EOD Api
	o := eod.Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-01-03"),
	}

	// Get a list of options
	options, underlyingLast, err := o.GetOptionsBySymbol("spy")

	// Double check results.
	st.Expect(t, err, nil)
	st.Expect(t, len(options), 4512)
	st.Expect(t, 270.47, underlyingLast)

	// Return a list of dates
	dates := bt.GetExpirationDatesFromOptions(options)

	// Test restults.
	st.Expect(t, len(dates), 30)
	st.Expect(t, dates[15].Format("2006-01-02"), "2018-04-20") // hehe 4/20 :)
}

//
// TestBackTestsSavingToDB01
//
func TestBackTestsSavingToDB01(t *testing.T) {
	// Only do this for non-short
	if testing.Short() {
		t.Skipf("Skipping TestDoPutCreditSpread01 test since it requires a env tokens and --short was requested")
	}

	// Load .env file (MUST CAll GO TEST FROM THE ROOT)
	env.ReadEnv("../.env")

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Build screener object
	screen := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 4.0},
			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.20},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// Set backtest
	btM := models.Backtest{
		UserId:          1,
		StartingBalance: 2000.00,
		EndingBalance:   2000.00,
		PositionSize:    "10-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2021-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2021-01-10")},
		Midpoint:        true,
		TradeSelect:     "least-days-to-expire",
		Screen:          screen,
	}

	// Run blank backtest
	bt := New(db, 1, "SPY")
	err := bt.DoBacktestDays(&btM)
	st.Expect(t, err, nil)
	st.Expect(t, len(btM.TradeGroups), 5)
	st.Expect(t, btM.CAGR, 743.9440649120804)
	st.Expect(t, btM.Return, 4.35)
	st.Expect(t, btM.Profit, 108.00)
	st.Expect(t, btM.BenchmarkCAGR, 285.21898330004444)
	st.Expect(t, btM.BenchmarkPercent, 3.3813281271184064)
	st.Expect(t, btM.BenchmarkEnd, 381.26)
	st.Expect(t, btM.BenchmarkStart, 368.79)

	// Save the backtest to database.
	db.Save(&btM)

	// Run again ane make sure the backtest is cleared.
	btM2 := models.Backtest{
		Id:              btM.Id,
		UserId:          1,
		StartingBalance: 2000.00,
		EndingBalance:   2000.00,
		PositionSize:    "10-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2021-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2021-01-10")},
		Midpoint:        true,
		TradeSelect:     "least-days-to-expire",
		Screen:          btM.Screen,
	}

	// Run blank backtest
	bt2 := New(db, 1, "SPY")
	err = bt2.DoBacktestDays(&btM2)
	st.Expect(t, err, nil)

	// Verify we saved correctly.
	var countBacktest int64
	var countBacktestPosition int64
	var countBacktestTradeGroup int64
	db.Model(&models.Backtest{}).Count(&countBacktest)
	db.Model(&models.BacktestPosition{}).Count(&countBacktestPosition)
	db.Model(&models.BacktestTradeGroup{}).Count(&countBacktestTradeGroup)
	st.Expect(t, countBacktest, int64(1))
	st.Expect(t, countBacktestPosition, int64(10))
	st.Expect(t, countBacktestTradeGroup, int64(5))
}

//
// TestRunPutCreditSpread01 - Run a put credit spread backtest. Testing "least-days-to-expire"
//
func TestDoPutCreditSpread01(t *testing.T) {
	// Only do this for non-short
	if testing.Short() {
		t.Skipf("Skipping TestDoPutCreditSpread01 test since it requires a env tokens and --short was requested")
	}

	// Load .env file (MUST CAll GO TEST FROM THE ROOT)
	env.ReadEnv("../.env")

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Build screener object
	screen := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 4.0},
			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.20},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// Set backtest
	btM := models.Backtest{
		UserId:          1,
		StartingBalance: 2000.00,
		EndingBalance:   2000.00,
		PositionSize:    "10-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2021-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2021-01-31")},
		Midpoint:        true,
		TradeSelect:     "least-days-to-expire",
		Benchmark:       "SPY",
		Screen:          screen,
	}

	// Run blank backtest
	bt := New(db, 1, btM.Benchmark)
	err := bt.DoBacktestDays(&btM)
	st.Expect(t, err, nil)
	st.Expect(t, len(btM.TradeGroups), 19)
	st.Expect(t, helpers.Round(btM.CAGR, 2), 842.69)
	st.Expect(t, btM.Return, 19.15)
	st.Expect(t, btM.Profit, 405.00)
	st.Expect(t, helpers.Round(btM.BenchmarkCAGR, 2), 4.31)
	st.Expect(t, helpers.Round(btM.BenchmarkPercent, 2), 0.35)
	st.Expect(t, btM.BenchmarkEnd, 370.07)
	st.Expect(t, btM.BenchmarkStart, 368.79)
}

//
// TestRunPutCreditSpread02 - Run a put credit spread backtest. Testing "highest-percent-away"
//
func TestDoPutCreditSpread02(t *testing.T) {
	// Only do this for non-short
	if testing.Short() {
		t.Skipf("Skipping TestDoPutCreditSpread02 test since it requires a env tokens and --short was requested")
	}

	// Load .env file (MUST CAll GO TEST FROM THE ROOT)
	env.ReadEnv("../.env")

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Build screener object
	screen := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 4.0},
			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.20},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// Set backtest
	btM := models.Backtest{
		UserId:          1,
		StartingBalance: 2000.00,
		EndingBalance:   2000.00,
		PositionSize:    "10-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2021-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2021-01-31")},
		Midpoint:        true,
		TradeSelect:     "highest-percent-away",
		Benchmark:       "SPY",
		Screen:          screen,
	}

	// Run blank backtest
	bt := New(db, 1, btM.Benchmark)
	err := bt.DoBacktestDays(&btM)
	st.Expect(t, err, nil)
	st.Expect(t, len(btM.TradeGroups), 19)
	st.Expect(t, helpers.Round(btM.CAGR, 2), 819.12)
	st.Expect(t, btM.Return, 18.90)
	st.Expect(t, btM.Profit, 400.00)
	st.Expect(t, helpers.Round(btM.BenchmarkCAGR, 2), 4.31)
	st.Expect(t, helpers.Round(btM.BenchmarkPercent, 2), 0.35)
	st.Expect(t, btM.BenchmarkEnd, 370.07)
	st.Expect(t, btM.BenchmarkStart, 368.79)
	st.Expect(t, btM.MaxDrawdown, 950.00)

	//bt.PrintResults(&btM)
}

//
// TestRunPutCreditSpread03 - Run a put credit spread backtest. Testing "shortest-percent-away"
//
func TestDoPutCreditSpread03(t *testing.T) {
	// Only do this for non-short
	if testing.Short() {
		t.Skipf("Skipping TestDoPutCreditSpread03 test since it requires a env tokens and --short was requested")
	}

	// Load .env file (MUST CAll GO TEST FROM THE ROOT)
	env.ReadEnv("../.env")

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Build screener object
	screen := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 4.0},
			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.20},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// Set backtest
	btM := models.Backtest{
		UserId:          1,
		StartingBalance: 5000.00,
		EndingBalance:   5000.00,
		PositionSize:    "10-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2020-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2020-04-01")},
		Midpoint:        true,
		TradeSelect:     "shortest-percent-away",
		Benchmark:       "SPY",
		Screen:          screen,
	}

	// Run blank backtest
	bt := New(db, 1, btM.Benchmark)
	err := bt.DoBacktestDays(&btM)
	st.Expect(t, err, nil)
	st.Expect(t, len(btM.TradeGroups), 51)
	st.Expect(t, helpers.Round(btM.CAGR, 2), 558.99)
	st.Expect(t, btM.Return, 58.54)
	st.Expect(t, btM.Profit, 2999.00)
	st.Expect(t, helpers.Round(btM.BenchmarkCAGR, 2), -67.15)
	st.Expect(t, helpers.Round(btM.BenchmarkPercent, 2), -24.22)
	st.Expect(t, btM.BenchmarkEnd, 246.15)
	st.Expect(t, btM.BenchmarkStart, 324.87)
	st.Expect(t, btM.MaxDrawdown, 950.00)

	bt.PrintResults(&btM)
}

// //
// // TestRunPutCreditSpread04 - Run a put credit spread backtest. Testing "max drawdown"
// //
// func TestDoPutCreditSpread04(t *testing.T) {
// 	// Only do this for non-short
// 	if testing.Short() {
// 		t.Skipf("Skipping TestDoPutCreditSpread04 test since it requires a env tokens and --short was requested")
// 	}

// 	// Load .env file (MUST CAll GO TEST FROM THE ROOT)
// 	env.ReadEnv("../.env")

// 	// Start the db connection.
// 	db, dbName, _ := models.NewTestDB("")
// 	defer models.TestingTearDown(db, dbName)

// 	// Build screener object
// 	screen := models.Screener{
// 		UserId:   1,
// 		Symbol:   "SPY",
// 		Strategy: "put-credit-spread",
// 		Items: []models.ScreenerItem{
// 			{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 4.0},
// 			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
// 			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
// 			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.20},
// 			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
// 			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
// 		},
// 	}

// 	// Set backtest
// 	btM := models.Backtest{
// 		UserId:          1,
// 		StartingBalance: 2000.00,
// 		EndingBalance:   2000.00,
// 		PositionSize:    "10-percent",
// 		StartDate:       models.Date{helpers.ParseDateNoError("2020-01-01")},
// 		EndDate:         models.Date{helpers.ParseDateNoError("2020-12-31")},
// 		Midpoint:        true,
// 		TradeSelect:     "least-days-to-expire",
// 		Benchmark:       "SPY",
// 		Screen:          screen,
// 	}

// 	// Run blank backtest
// 	bt := New(db, 1, btM.Benchmark)
// 	err := bt.DoBacktestDays(&btM)
// 	st.Expect(t, err, nil)
// 	st.Expect(t, len(btM.TradeGroups), 19)
// 	// st.Expect(t, helpers.Round(btM.CAGR, 2), 842.69)
// 	// st.Expect(t, btM.Return, 19.15)
// 	// st.Expect(t, btM.Profit, 405.00)
// 	// st.Expect(t, helpers.Round(btM.BenchmarkCAGR, 2), 4.31)
// 	// st.Expect(t, helpers.Round(btM.BenchmarkPercent, 2), 0.35)
// 	// st.Expect(t, btM.BenchmarkEnd, 370.07)
// 	// st.Expect(t, btM.BenchmarkStart, 368.79)
// 	st.Expect(t, btM.MaxDrawdown, 60.79)
// }

/* End File */
