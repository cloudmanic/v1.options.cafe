//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"fmt"
	"time"

	"bitbucket.org/api.triwou.org/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/screener"
)

const workerCount int = 10

// Used to cache all symbols in the DB to avoid too many sql queries
var cachedSymbols map[string]models.Symbol = make(map[string]models.Symbol)

// Base struct
type Base struct {
	DB           models.Datastore
	ResultsFuncs map[string]func(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) ([]screener.Result, error)
}

// Job struct
type Job struct {
	Day      time.Time
	Index    int
	Backtest *models.Backtest
	Results  []screener.Result
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
	t.ResultsFuncs = map[string]func(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) ([]screener.Result, error){
		"blank":             t.DoBlank,
		"put-credit-spread": t.DoPutCreditSpread,
	}

	// Warm symbol cache (just random sumbol to warm cache)
	t.GetSymbol("SPY190418P00269000", "SPY Apr 18 2019 $269.00 Put", "Option")

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

	// Worker job queue (some stocks have thousands of days)
	totalJobs := 0
	jobs := make(chan Job, 1000000)
	results := make(chan Job, 1000000)

	// Load up the workers
	for w := 0; w < workerCount; w++ {
		go t.tradeResultsWorker(jobs, results)
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

		//go func() {

		jobs <- Job{Index: totalJobs, Day: row, Backtest: backtest}
		totalJobs++

		// // Job struct
		// type Job struct {
		// 	Symbol         string
		// 	Day            time.Time
		// 	Index          int
		// 	UnderlyingLast float64
		// 	Options        []types.OptionsChainItem
		// }

		// // Create broker object
		// o := eod.Api{
		// 	DB:  t.DB,
		// 	Day: row,
		// }
		//
		// // Log where we are in the backtest. TODO(spicer): Send this up websocket
		// services.Info("Backtesting " + backtest.Screen.Strategy + " " + backtest.Screen.Symbol + " on " + row.Format("2006-01-02"))
		//
		// // Get all options for this symbol and day.
		// options, underlyingLast, _ := o.GetOptionsBySymbol(backtest.Screen.Symbol)
		//
		// // if err != nil {
		// // 	return err
		// // }
		//
		// // Run backtest strategy function for this backtest
		// t.StrategyFuncs[backtest.Screen.Strategy](row, backtest, underlyingLast, options)
		//
		// // if err != nil {
		// // 	return err
		// // }

		//}()

	}

	// Close jobs so the workers return.
	close(jobs)

	// Set the results array list.
	resultsList := make([][]screener.Result, totalJobs)

	// Collect results so this function does not just return.
	for a := 0; a < totalJobs; a++ {
		job := <-results
		resultsList[job.Index] = job.Results
	}

	// Close results
	close(results)

	// Now that we have a list of all possible trades by day we can go through and "trade".
	for _, row := range resultsList {
		for _, row2 := range row {
			fmt.Println(row2.Credit)
		}
	}

	// Store how long the backtest took to run.
	backtest.TimeElapsed = time.Since(start)

	return nil
}

//
// tradeResultsWorker - A worker for running trade results per day.
//
func (t *Base) tradeResultsWorker(jobs <-chan Job, results chan<- Job) {
	// Wait for jobs to come in and process them.
	for job := range jobs {
		// Create broker object
		o := eod.Api{
			DB:  t.DB,
			Day: job.Day,
		}

		// Log where we are in the backtest. TODO(spicer): Send this up websocket
		services.Info("Backtesting " + job.Backtest.Screen.Strategy + " " + job.Backtest.Screen.Symbol + " on " + job.Day.Format("2006-01-02"))

		// Get all options for this symbol and day.
		options, underlyingLast, err := o.GetOptionsBySymbol(job.Backtest.Screen.Symbol)

		if err != nil {
			services.Warning(err)
		}

		// Run backtest strategy function for this backtest
		job.Results, err = t.ResultsFuncs[job.Backtest.Screen.Strategy](job.Day, job.Backtest, underlyingLast, options)

		if err != nil {
			services.Warning(err)
		}

		// Send back a happy with results.
		results <- job
	}
}

//
// DoBlank is mostly using for unit testing.
//
func (t *Base) DoBlank(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) ([]screener.Result, error) {
	return []screener.Result{}, nil
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