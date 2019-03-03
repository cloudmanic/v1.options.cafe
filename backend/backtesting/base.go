//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

const workerCount int = 100

// Used to cache all symbols in the DB to avoid too many sql queries
var cachedSymbols map[string]models.Symbol = make(map[string]models.Symbol)

// Base struct
type Base struct {
	DB            models.Datastore
	StrategyFuncs map[string]func(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) error
}

// Job struct
type Job struct {
	Symbol         string
	Day            time.Time
	Index          int
	UnderlyingLast float64
	Options        []types.OptionsChainItem
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
	// Let's time this backtest
	start := time.Now()

	// Get dates by symbol
	dates, err := eod.GetTradeDatesBySymbols(backtest.Screen.Symbol)

	if err != nil {
		return err
	}

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
		//start2 := time.Now()
		options, underlyingLast, err := o.GetOptionsBySymbol(backtest.Screen.Symbol)
		//elapsed2 := time.Since(start2)
		//log.Printf("BS took %s", elapsed2)
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

	// Store how long the backtest took to run.
	backtest.TimeElapsed = time.Since(start)
	//log.Printf("Binomial took %s", backtest.TimeElapsed)
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

//
// GetSymbol - This is a wrapper function for models.Symbol. We want to do a bit of "caching".
//
func (t *Base) GetSymbol(short string, name string, sType string) (models.Symbol, error) {
	// Build cache of all symbols in the system
	if len(cachedSymbols) == 0 {
		// Get all the symbols in the DB
		s := t.DB.GetAllSymbols()

		// Loop through and build hash table
		for _, row := range s {
			cachedSymbols[row.ShortName] = row
		}
	}

	// See if we have this symbol in cache. If so return happy.
	if val, ok := cachedSymbols[short]; ok {
		return val, nil
	}

	// Add symbol to the DB. Since we do not know about it.
	symb, err := t.DB.CreateNewSymbol(short, name, sType)

	if err != nil {
		return symb, err
	}

	// Add symbol to map.
	cachedSymbols[short] = symb

	// Return happy.
	return symb, nil
}

/* End File */
