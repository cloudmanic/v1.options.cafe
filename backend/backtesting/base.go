//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"time"

	longstraddle "app.options.cafe/backtesting/long-straddle"
	"app.options.cafe/brokers/eod"
	"app.options.cafe/brokers/tradier"
	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/queue"
	"app.options.cafe/library/services"
	"app.options.cafe/library/worker"
	"app.options.cafe/models"
	"app.options.cafe/screener"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"

	screenerCache "app.options.cafe/screener/cache"
)

const workerCount int = 3

// Base struct
type Base struct {
	DB              models.Datastore
	UserID          int
	BenchmarkQuotes []types.HistoryQuote
	TradeFuncs      map[string]func(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem)
	ResultsFuncs    map[string]func(db models.Datastore, today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error)
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
// BacktestDaysWorker will run a backtest days job from the worker queue
//
// 	queue.Write("oc-job", `{"action":"backtest-run-days","user_id":`+strconv.Itoa(userID)+`,"backtest_id":`+strconv.Itoa(77)+`}`)
//
func BacktestDaysWorker(job worker.JobRequest) error {
	// Log to console.
	services.InfoMsg("Starting Backtest " + strconv.Itoa(int(job.BacktestId)))

	// Get the backtest
	btM, err := job.DB.BacktestGetById(job.BacktestId)

	if err != nil {
		return err
	}

	// Setup a new backtesting & run it.
	bt := New(job.DB, int(btM.UserId), btM.Benchmark)
	bt.DoBacktestDays(&btM)

	// Log to console.
	services.InfoMsg("Ending Backtest " + strconv.Itoa(int(job.BacktestId)))

	return nil
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
	t.ResultsFuncs = map[string]func(db models.Datastore, today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error){
		"empty":                      t.EmptyResults, // used for testing
		"long-straddle":              longstraddle.Results,
		"put-credit-spread":          t.PutCreditSpreadResults,
		"long-call-butterfly-spread": t.LongCallButterflySpreadResults,
	}

	// Build backtest functions - trades
	t.TradeFuncs = map[string]func(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem){
		"empty":                      t.EmptyTrades, // used for testing
		"long-straddle":              longstraddle.Trades,
		"put-credit-spread":          t.PutCreditSpreadPlaceTrades,
		"long-call-butterfly-spread": t.LongCallButterflySpreadPlaceTrades,
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

	// Clear past backtests
	t.ClearPastBacktests(backtest)

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

	// CAGR = (Ending value / Starting value)^(1 / # years) -1
	s1 := helpers.ParseDateNoError(backtest.StartDate.Format("01/02/2006"))
	s2 := helpers.ParseDateNoError(backtest.EndDate.Format("01/02/2006"))
	years := (s2.Sub(s1).Hours() / 24) / 365
	backtest.CAGR = (math.Pow((backtest.EndingBalance/backtest.StartingBalance), (1/years)) - 1) * 100
	backtest.BenchmarkCAGR = (math.Pow((backtest.BenchmarkEnd/backtest.BenchmarkStart), (1/years)) - 1) * 100

	// Returns / win ratio
	wins := 0
	backtest.Return = 0.00

	for _, row := range backtest.TradeGroups {
		// Only record closed position
		if row.Status == "Closed" {
			backtest.Return = row.ReturnFromStart

			if row.ReturnPercent > 0 {
				wins++
			}
		}
	}

	// Backtest stats
	backtest.Profit = (backtest.EndingBalance - backtest.StartingBalance)
	backtest.TradeCount = len(backtest.TradeGroups)
	backtest.BenchmarkPercent = (((backtest.BenchmarkEnd - backtest.BenchmarkStart) / backtest.BenchmarkStart) * 100)
	backtest.WinRatio = (float64(wins) / float64(backtest.TradeCount)) * 100

	// Store how long the backtest took to run.
	backtest.TimeElapsed = time.Since(start)

	// Save backtest to DB. If we run as a CMD we might not save
	if backtest.Id > 0 {
		t.DB.New().Save(backtest)
	}

	// Display results. Just used for debugging.
	t.PrintResults(backtest)

	return nil
}

//
// ClearPastBacktests will clear data for past backtests.
//
func (t *Base) ClearPastBacktests(backtest *models.Backtest) {
	// Maybe we are rerunning the backtest so we clear out past data
	if backtest.Id > 0 {
		pIds := []uint{}

		bTmp := models.Backtest{Id: backtest.Id}
		t.DB.New().Preload("TradeGroups").First(&bTmp)

		for _, row := range bTmp.TradeGroups {
			pIds = append(pIds, row.Id)
		}

		t.DB.New().Exec("DELETE FROM backtest_trade_groups WHERE backtest_id = ?", backtest.Id)
		t.DB.New().Exec("DELETE FROM backtest_positions WHERE backtest_trade_group_id IN (?)", pIds)

		backtest.Profit = 0.00
		backtest.Return = 0.00
		backtest.CAGR = 0.00
		backtest.TradeCount = 0
		backtest.WinRatio = 0.00
		backtest.TimeElapsed = 0
		backtest.BenchmarkStart = 0.00
		backtest.BenchmarkEnd = 0.00
		backtest.BenchmarkCAGR = 0.00
		backtest.EndingBalance = backtest.StartingBalance // No trades means starting and ending are the same
		backtest.BenchmarkPercent = 0.00
		backtest.TradeGroups = []models.BacktestTradeGroup{}

		t.DB.New().Save(backtest)
	}
}

//
// PrintResults will print results to the screen. Useful for debugging and running from CLI
//
func (t *Base) PrintResults(backtest *models.Backtest) error {
	plotData := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	csvData := [][]string{{"Open Date", "Close Date", "Spread", "Open", "Close", "Credit", "Return", "Lots", "% Away", "Margin", "Balance", "Return", "Benchmark", "Benchmark Balance", "Benchmark Return", "Status", "Note"}}
	table.SetHeader([]string{"Open Date", "Close Date", "Spread", "Open", "Close", "Credit", "Return", "Lots", "% Away", "Margin", "Balance", "Return", "Benchmark", "Benchmark Balance", "Benchmark Return", "Status", "Note"})

	// Build position rows
	for _, row := range backtest.TradeGroups {
		d := []string{
			row.OpenDate.Format("01/02/2006"),
			row.CloseDate.Format("01/02/2006"),
			fmt.Sprintf("%s %s %.2f / %.2f", row.Positions[0].Symbol.OptionUnderlying, row.Positions[0].Symbol.OptionExpire.Format("01/02/2006"), row.Positions[0].Symbol.OptionStrike, row.Positions[1].Symbol.OptionStrike),
			fmt.Sprintf("$%.2f", row.OpenPrice),
			fmt.Sprintf("$%.2f", row.ClosePrice),
			fmt.Sprintf("$%.2f", row.Credit),
			fmt.Sprintf("%.2f", row.ReturnPercent) + "%",
			fmt.Sprintf("%d", row.Lots),
			fmt.Sprintf("%.2f", row.PutPrecentAway) + "%",
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(row.Margin))),
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(row.Balance))),
			fmt.Sprintf("%.2f", row.ReturnFromStart) + "%",
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(row.BenchmarkLast))),
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(row.BenchmarkBalance))),
			fmt.Sprintf("%.2f", row.BenchmarkReturn) + "%",
			row.Status,
			row.Note,
		}

		table.Append(d)
		csvData = append(csvData, d)
		plotData = append(plotData, []string{row.OpenDate.Format("01/02/2006"), humanize.BigCommaf(big.NewFloat(row.Balance))})
	}
	table.Render()

	// Show how long the backtest took.
	log.Printf("Backtest took %s", backtest.TimeElapsed)
	log.Println("")
	log.Println("Summmary")
	log.Println("-------------")
	log.Printf("CAGR: %s%%", humanize.BigCommaf(big.NewFloat(helpers.Round(backtest.CAGR, 2))))
	log.Printf("Return: %s%%", humanize.BigCommaf(big.NewFloat(helpers.Round(backtest.Return, 2))))
	log.Printf("Profit: $%s", humanize.BigCommaf(big.NewFloat(helpers.Round(backtest.Profit, 2))))
	log.Printf("Trade Count: %d", backtest.TradeCount)
	log.Println("")
	log.Println("Benchmark")
	log.Println("-------------")
	log.Printf("Start (%s): %s", backtest.Benchmark, humanize.BigCommaf(big.NewFloat(helpers.Round(backtest.BenchmarkStart, 2))))
	log.Printf("End (%s): %s", backtest.Benchmark, humanize.BigCommaf(big.NewFloat(helpers.Round(backtest.BenchmarkEnd, 2))))
	log.Printf("CAGR (%s): %s%%", backtest.Benchmark, humanize.BigCommaf(big.NewFloat(helpers.Round(backtest.BenchmarkCAGR, 2))))
	log.Printf("Return (%s): %s%%", backtest.Benchmark, humanize.BigCommaf(big.NewFloat(helpers.Round(backtest.BenchmarkPercent, 2))))
	log.Println("")

	return nil
}

//
// GetOptionsByExpirationType - Loop through and filter out just expire and type. TODO(spicer): Kill. This was moved to helpers.
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
// GetExpirationDatesFromOptions - Take complete list of options and return a list of expiration dates. TODO(spicer): Kill. This was moved to helpers.
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
		job.Results, err = t.ResultsFuncs[job.Backtest.Screen.Strategy](t.DB, job.Day, job.Backtest, underlyingLast, options, cache)

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
