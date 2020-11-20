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
)

//
// TestDoBacktestDays01 - Run a backtest looping through each day
//
func TestDoBacktestDays01(t *testing.T) {

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Setup a new backtesting
	bt := New(db, "SPY")

	// Build screener object
	screen := models.Screener{
		Symbol:   "SPY",
		Strategy: "blank",
	}

	// Run blank backtest
	err := bt.DoBacktestDays(&models.Backtest{
		StartDate: models.Date{helpers.ParseDateNoError("2018-01-01")},
		EndDate:   models.Date{helpers.ParseDateNoError("2019-01-01")},
		Screen:    screen,
	})
	st.Expect(t, err, nil)
}

//
// TestGetOptionsByExpirationType01
//
func TestGetOptionsByExpirationType01(t *testing.T) {
	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Setup a new backtesting
	bt := New(db, "SPY")

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
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Setup a new backtesting
	bt := New(db, "SPY")

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
// TestGetSymbol - Get getting symb struct from DB
//
func TestGetSymbol(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Setup a new backtesting
	bt := New(db, "SPY")

	// Get test symbol.
	smb, err := bt.GetSymbol("SPY190418P00269000", "SPY Apr 18 2019 $269.00 Put", "Option")

	// Check results.
	st.Expect(t, err, nil)
	st.Expect(t, smb.ShortName, "SPY190418P00269000")
	st.Expect(t, smb.Name, "SPY Apr 18 2019 $269.00 Put")
	st.Expect(t, smb.OptionStrike, 269.00)

	// Get test symbol - Make sure it is cached
	smb2, err := bt.GetSymbol("SPY190418P00269000", "SPY Apr 18 2019 $269.00 Put", "Option")

	// Check results.
	st.Expect(t, err, nil)
	st.Expect(t, smb2.Id, smb.Id)
	st.Expect(t, smb2.ShortName, "SPY190418P00269000")
	st.Expect(t, smb2.Name, "SPY Apr 18 2019 $269.00 Put")
	st.Expect(t, smb2.OptionStrike, 269.00)
}

/* End File */
