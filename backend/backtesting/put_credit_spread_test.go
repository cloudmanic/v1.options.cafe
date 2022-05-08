//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"testing"

	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
	"github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
)

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
	st.Expect(t, btM.CAGR, 847.4692690693325)
	st.Expect(t, btM.Return, 13.80)
	st.Expect(t, btM.Profit, 406.00)
	st.Expect(t, btM.BenchmarkCAGR, 4.305621689814987)
	st.Expect(t, btM.BenchmarkPercent, 0.34708099460396774)
	st.Expect(t, btM.BenchmarkEnd, 370.07)
	st.Expect(t, btM.BenchmarkStart, 368.79)
}

//
// TestRunPutCreditSpread02 - Run a put credit spread backtest. Testing "highest-percent-away"
//
func TestDoPutCreditSpread02(t *testing.T) {
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
		TradeSelect:     "highest-percent-away",
		Benchmark:       "SPY",
		Screen:          screen,
	}

	// Run blank backtest
	bt := New(db, 1, btM.Benchmark)
	err := bt.DoBacktestDays(&btM)
	st.Expect(t, err, nil)
	st.Expect(t, len(btM.TradeGroups), 12)
	st.Expect(t, btM.CAGR, 332.9396672739466)
	st.Expect(t, btM.Return, float64(0))
	st.Expect(t, btM.Profit, 256.00)
	st.Expect(t, btM.BenchmarkCAGR, 4.305621689814987)
	st.Expect(t, btM.BenchmarkPercent, 0.34708099460396774)
	st.Expect(t, btM.BenchmarkEnd, 370.07)
	st.Expect(t, btM.BenchmarkStart, 368.79)
}

//
// TestRunPutCreditSpread03 - Run a put credit spread backtest. Testing "shortest-percent-away"
//
func TestDoPutCreditSpread03(t *testing.T) {
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
		StartDate:       models.Date{helpers.ParseDateNoError("2017-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2021-12-31")},
		Midpoint:        true,
		TradeSelect:     "shortest-percent-away",
		Benchmark:       "SPY",
		Screen:          screen,
	}

	// Run blank backtest
	bt := New(db, 1, btM.Benchmark)
	err := bt.DoBacktestDays(&btM)
	st.Expect(t, err, nil)
	st.Expect(t, len(btM.TradeGroups), 19)
	st.Expect(t, btM.CAGR, 937.4054883212459)
	st.Expect(t, btM.Return, 15.35)
	st.Expect(t, btM.Profit, 424.00)
	st.Expect(t, btM.BenchmarkCAGR, 4.305621689814987)
	st.Expect(t, btM.BenchmarkPercent, 0.34708099460396774)
	st.Expect(t, btM.BenchmarkEnd, 370.07)
	st.Expect(t, btM.BenchmarkStart, 368.79)
}

/* End File */
