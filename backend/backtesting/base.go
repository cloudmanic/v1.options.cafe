//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"os"
	"strconv"
	"time"

	"app.options.cafe/brokers/eod"
	"app.options.cafe/brokers/tradier"
	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/queue"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
	"app.options.cafe/screener"

	screenerCache "app.options.cafe/screener/cache"
)

const workerCount int = 3

// Base struct
type Base struct {
	DB              models.Datastore
	UserID          int
	BenchmarkQuotes []types.HistoryQuote
	TradeFuncs      map[string]func(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem)
	ResultsFuncs    map[string]func(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error)
}

// Job struct
type Job struct {
	Day       time.Time
	Index     int
	TotalDays int
	Backtest  *models.Backtest
	Results   []screener.Result
	Options   []types.OptionsChainItem
}

//
// New Backtest
//
func New(db models.Datastore, userID int, benchmark string) Base {
	// New backtest instance
	t := Base{
		DB:     db,
		UserID: userID,
	}

	// Build backtest functions - results
	t.ResultsFuncs = map[string]func(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error){
		"put-credit-spread": t.PutCreditSpreadResults,
	}

	// Build backtest functions - trades
	t.TradeFuncs = map[string]func(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem){
		"put-credit-spread": t.PutCreditSpreadPlaceTrades,
	}

	// Get benchmark data.
	tr := &tradier.Api{ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")}
	startDate := time.Date(2009, 01, 01, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(3000, 12, 31, 0, 0, 0, 0, time.UTC) // Some big future date, I hope someday this is a bug :)
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

	// Send up websocket
	queue.Write("oc-websocket-write", `{"uri":"backtest-start","user_id":`+strconv.Itoa(t.UserID)+`,"body":`+helpers.JsonEncode(backtest.Screen)+`}`)

	// Build the cache for this screen.
	cache := screenerCache.New(t.DB)

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
		go t.tradeResultsWorker(jobs, results, cache)
	}

	// Loop through and figure out what dates we need to test.
	days := []time.Time{}

	for _, row := range dates {
		// We skip dates before our start date
		if row.Before(helpers.ParseDateNoError(backtest.StartDate.Format("2006-01-02"))) {
			continue
		}

		// We skip dates after our end date
		if (row.Format("2006-01-02") != backtest.EndDate.Format("2006-01-02")) && row.After(helpers.ParseDateNoError(backtest.EndDate.Format("2006-01-02"))) {
			continue
		}

		days = append(days, row)
		totalJobs++
	}

	// Loop through the dates and run backtest.
	c := 0

	for _, row := range days {
		// Set the benchmark start
		backtest.BenchmarkEnd = t.getBenchmarkByDate(helpers.ParseDateNoError(row.Format("2006-01-02")))

		// Figure out benchmark start value
		if backtest.BenchmarkStart == 0.00 {
			backtest.BenchmarkStart = t.getBenchmarkByDate(helpers.ParseDateNoError(row.Format("2006-01-02")))
		}

		// Add job to worker queue
		jobs <- Job{Index: c, TotalDays: totalJobs, Day: row, Backtest: backtest}

		// Update index
		c++
	}

	// Send total days to websocket
	queue.Write("oc-websocket-write", `{"uri":"backtest-total-days","user_id":`+strconv.Itoa(t.UserID)+`,"body":{"backtest_id":`+strconv.Itoa(int(backtest.Id))+`,"total_days":`+strconv.Itoa(totalJobs)+`}}`)

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

	// Send up websocket
	queue.Write("oc-websocket-write", `{"uri":"backtest-end","user_id":`+strconv.Itoa(t.UserID)+`,"body":`+helpers.JsonEncode(backtest.Screen)+`}`)

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
func (t *Base) tradeResultsWorker(jobs <-chan Job, results chan<- Job, cache screenerCache.Cache) {
	// Wait for jobs to come in and process them.
	for job := range jobs {
		// Create broker object
		o := eod.Api{
			DB:  t.DB,
			Day: job.Day,
		}

		// Websocket MSG
		type Msg struct {
			Day       time.Time
			Screen    models.Screener
			TotalDays int
			Index     int
		}

		msg := Msg{
			Day:       job.Day,
			Screen:    job.Backtest.Screen,
			Index:     job.Index,
			TotalDays: job.TotalDays,
		}

		// Send up websocket
		queue.Write("oc-websocket-write", `{"uri":"backtest-day-run","user_id":`+strconv.Itoa(t.UserID)+`,"body":`+helpers.JsonEncode(msg)+`}`)

		// Log where we are in the backtest.
		services.InfoMsg("Backtesting " + job.Backtest.Screen.Strategy + " " + job.Backtest.Screen.Symbol + " on " + job.Day.Format("2006-01-02") + " (" + strconv.Itoa(job.Index) + "/" + strconv.Itoa(job.TotalDays) + ")")

		// Get all options for this symbol and day.
		options, underlyingLast, err := o.GetOptionsBySymbol(job.Backtest.Screen.Symbol)

		if err != nil {
			services.Info(err)
		}

		// Run backtest strategy function for this backtest
		job.Results, err = t.ResultsFuncs[job.Backtest.Screen.Strategy](job.Day, job.Backtest, underlyingLast, options, cache)

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
