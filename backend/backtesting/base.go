//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"os"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/screener"
)

const workerCount int = 3

// Base struct
type Base struct {
	DB              models.Datastore
	BenchmarkQuotes []types.HistoryQuote
	CachedSymbols   map[string]models.Symbol
	TradeFuncs      map[string]func(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem)
	ResultsFuncs    map[string]func(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) ([]screener.Result, error)
}

// Job struct
type Job struct {
	Day      time.Time
	Index    int
	Backtest *models.Backtest
	Results  []screener.Result
	Options  []types.OptionsChainItem
}

//
// New Backtest
//
func New(db models.Datastore, benchmark string) Base {
	// New backtest instance
	t := Base{
		DB: db,
	}

	// Build cache of all symbols in the system
	s := t.DB.GetAllSymbols()
	t.CachedSymbols = make(map[string]models.Symbol)

	// Loop through and build hash table
	for _, row := range s {
		t.CachedSymbols[row.ShortName] = row
	}

	// Build backtest functions - results
	t.ResultsFuncs = map[string]func(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) ([]screener.Result, error){
		"put-credit-spread": t.PutCreditSpreadResults,
	}

	// Build backtest functions - trades
	t.TradeFuncs = map[string]func(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem){
		"put-credit-spread": t.PutCreditSpreadPlaceTrades,
	}

	// Warm symbol cache (just random sumbol to warm cache)
	t.GetSymbol("SPY190418P00269000", "SPY Apr 18 2019 $269.00 Put", "Option")

	// Get benchmark data.
	tr := &tradier.Api{ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")}
	startDate := time.Date(2009, 01, 01, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	t.BenchmarkQuotes, _ = tr.GetHistoricalQuotes(benchmark, startDate, endDate, "daily")

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

		// Add job to worker queue
		jobs <- Job{Index: totalJobs, Day: row, Backtest: backtest}
		totalJobs++
	}

	// Close jobs so the workers return.
	close(jobs)

	// Set the results array list.
	resultsList := make([]Job, totalJobs)

	// Collect results so this function does not just return.
	for a := 0; a < totalJobs; a++ {
		job := <-results
		resultsList[job.Index] = job
	}

	// Close results
	close(results)

	// Now that we have a list of all possible trades by day we can go through and "trade".
	for _, row := range resultsList {
		// Send the results into a trade function.
		t.TradeFuncs[backtest.Screen.Strategy](row.Day, backtest, row.Results, row.Options)
	}

	// Store how long the backtest took to run.
	backtest.TimeElapsed = time.Since(start)

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
	// See if we have this symbol in cache. If so return happy.
	if val, ok := t.CachedSymbols[short]; ok {
		return val, nil
	}

	// Add symbol to the DB. Since we do not know about it.
	symb, err := t.DB.CreateNewSymbol(short, name, sType)

	if err != nil {
		return symb, err
	}

	// Add symbol to map. TODO(spicer) have to do magic to manage concurrent access with go routines
	//t.CachedSymbols[short] = symb

	// Return happy.
	return symb, nil
}

// ---------------- Private helper functions ------------------ //

//
// getBenchmarkByDate will return the close price of the benchmark
//
func (t *Base) getBenchmarkByDate(date time.Time) float64 {
	for _, row := range t.BenchmarkQuotes {
		if row.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			return row.Close
		}
	}
	return 0.00
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
		services.InfoMsg("Backtesting " + job.Backtest.Screen.Strategy + " " + job.Backtest.Screen.Symbol + " on " + job.Day.Format("2006-01-02"))

		// Get all options for this symbol and day.
		options, underlyingLast, err := o.GetOptionsBySymbol(job.Backtest.Screen.Symbol)

		if err != nil {
			services.Info(err)
		}

		// Run backtest strategy function for this backtest
		job.Results, err = t.ResultsFuncs[job.Backtest.Screen.Strategy](job.Day, job.Backtest, underlyingLast, options)

		if err != nil {
			services.Info(err)
		}

		// Add options to job
		job.Options = options

		// Send back a happy with results.
		results <- job
	}
}

/* End File */
