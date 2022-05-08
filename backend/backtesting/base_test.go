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

/* End File */
