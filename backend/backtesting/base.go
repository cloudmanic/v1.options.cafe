//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"log"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

// Base struct
type Base struct {
	DB            models.Datastore
	StrategyFuncs map[string]func(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) error
}

//
// New Backtest
//
func New(db models.Datastore) Base {

	// New backtest instance
	t := Base{
		DB: db,
	}

	// Build backtest functions.
	t.StrategyFuncs = map[string]func(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) error{
		"blank":             t.DoBlank,
		"put-credit-spread": t.DoPutCreditSpread,
	}

	return t
}

//
// DoBacktestDays - Loop through every day in the backtest and pass
// an options chain into a function for a prticular backtest type.
//
func (t *Base) DoBacktestDays(backtest *models.Backtest) error {

	// Get dates by symbol
	dates, err := eod.GetTradeDatesBySymbols(backtest.Screen.Symbol)

	if err != nil {
		return err
	}

	start := time.Now()

	// Loop through the dates and run backtest.
	for _, row := range dates {

		// We skip dates before our start date
		if row.Before(helpers.ParseDateNoError(backtest.StartDate.Format("2006-01-02"))) {
			continue
		}

		// We skip dates after our end date
		if row.After(helpers.ParseDateNoError(backtest.EndDate.Format("2006-01-02"))) {
			continue
		}

		// Create broker object
		o := eod.Api{
			DB:  t.DB,
			Day: row,
		}

		// Log where we are in the backtest. TODO(spicer): Send this up websocket
		services.Info("Backtesting " + backtest.Screen.Strategy + " " + backtest.Screen.Symbol + " on " + row.Format("2006-01-02"))

		// Get all options for this symbol and day.
		//start := time.Now()
		options, underlyingLast, err := o.GetOptionsBySymbol(backtest.Screen.Symbol)
		// elapsed := time.Since(start)
		// log.Printf("Binomial took %s", elapsed)
		// os.Exit(1)

		if err != nil {
			return err
		}

		// Run backtest strategy function for this backtest
		err = t.StrategyFuncs[backtest.Screen.Strategy](row, backtest, underlyingLast, options)

		if err != nil {
			return err
		}
	}

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
	//os.Exit(1)

	return nil
}

//
// DoBlank is mostly using for unit testing.
//
func (t *Base) DoBlank(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) error {
	return nil
}

//
// GetOptionsByExpirationType - Loop through and filter out just expire and type
//
func (t *Base) GetOptionsByExpirationType(expire types.Date, optionType string, options []types.OptionsChainItem) []types.OptionsChainItem {
	rt := []types.OptionsChainItem{}

	// Double check TODO(spicer): Return error maybe
	if (optionType != "Put") && (optionType != "Call") {
		return rt
	}

	for _, row := range options {

		if row.OptionType != optionType {
			continue
		}

		if row.ExpirationDate.Format("2006-01-02") != expire.Format("2006-01-02") {
			continue
		}

		rt = append(rt, row)
	}

	// Return filtered subset
	return rt
}

//
// GetExpirationDatesFromOptions - Take complete list of options and return a list of expiration dates.
//
func (t *Base) GetExpirationDatesFromOptions(options []types.OptionsChainItem) []types.Date {
	seen := map[types.Date]bool{}
	dates := []types.Date{}

	// Loop through and get expire dates
	for _, row := range options {
		if _, ok := seen[row.ExpirationDate]; !ok {
			dates = append(dates, row.ExpirationDate)
			seen[row.ExpirationDate] = true
		}
	}

	// Return happy
	return dates
}

/* End File */
