//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// TestDoBacktestDays01 - Run a backtest looping through each day
//
func TestDoBacktestDays01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Setup a new backtesting
	bt := New(db)

	// Build screener object
	screen := models.Screener{
		Symbol:   "SPY",
		Strategy: "blank",
	}

	// Run blank backtest
	err := bt.DoBacktestDays(&models.Backtest{
		StartDate: models.Date{helpers.ParseDateNoError("2018-01-01")},
		EndDate:   models.Date{helpers.ParseDateNoError("2018-01-03")},
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
	bt := New(db)

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

/* End File */
