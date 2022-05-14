//
// Date: 2019-02-27
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package actions

import (
	"app.options.cafe/backtesting"
	"app.options.cafe/library/worker"
	"app.options.cafe/models"
)

//
// RunBackTest run a back test
//
// go run main.go --cmd="backtest-run" --user_id=1
//
// gnuplot -e "set terminal jpeg size 3600,1600; set timefmt '%m/%d/%Y'; set datafile separator ','; set format x '%m/%d/%Y'; set xdata time; set title 'Account Balance Over Time';  plot for [col=2:2] '/Users/spicer/Downloads/graph.csv' using 1:col with lines" > /tmp/blah.jpeg && open /tmp/blah.jpeg
//
func RunBackTest(db *models.DB, userID int) {

	// Send directly to the worker without a queue.
	backtesting.BacktestDaysWorker(worker.JobRequest{DB: db, BacktestId: 1})

	// Send to the worker queue
	//queue.Write("oc-job", `{"action":"backtest-run-days","backtest_id":`+strconv.Itoa(3)+`}`)

	// // Build screener object
	// screen := models.Screener{
	// 	UserId:   1,
	// 	Symbol:   "SPY",
	// 	Name: "SPY Percent Away 45 Days",
	// 	Strategy: "put-credit-spread",
	// 	Items: []models.ScreenerItem{
	// 		{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 4.5},
	// 		{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
	// 		{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
	// 		{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.50},
	// 		{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
	// 		{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
	// 		{UserId: 1, Key: "allow-more-than-one-expire", Operator: "=", ValueString: "no"},
	// 		{UserId: 1, Key: "allow-more-than-one-strike", Operator: "=", ValueString: "no"},
	// 	},
	// }

	// // Set backtest
	// btM := models.Backtest{
	// 	UserId:          uint(userID),
	// 	StartingBalance: 5000.00,
	// 	EndingBalance:   5000.00,
	// 	PositionSize:    "15-percent", // one-at-time, *-percent
	// 	StartDate:       models.Date{helpers.ParseDateNoError("2022-01-01")},
	// 	EndDate:         models.Date{helpers.ParseDateNoError("2022-12-31")},
	// 	Midpoint:        true,
	// 	TradeSelect:     "least-days-to-expire", // least-days-to-expire, highest-midpoint, highest-ask, highest-percent-away, shortest-percent-away
	// 	Benchmark:       "SPY",
	// 	Screen:          screen,
	// }

	// // Setup a new backtesting
	// bt := backtesting.New(db, userID, btM.Benchmark)

	// // Run the backtest
	// bt.DoBacktestDays(&btM)

	// // Display results
	// bt.PrintResults(&btM)
}

/* End File */
