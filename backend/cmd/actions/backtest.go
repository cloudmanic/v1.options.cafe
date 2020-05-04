//
// Date: 2019-02-27
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package actions

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/backtesting"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/olekukonko/tablewriter"
	"github.com/optionscafe/options-cafe-cli/helpers"

	humanize "github.com/dustin/go-humanize"
)

//
// RunBackTest run a back test
//
// go run main.go --cmd="backtest-run" --user_id=1
//
// gnuplot -e "set terminal jpeg size 3600,1600; set timefmt '%m/%d/%Y'; set datafile separator ','; set format x '%m/%d/%Y'; set xdata time; set title 'Account Balance Over Time';  plot for [col=2:2] '/Users/spicer/Downloads/graph.csv' using 1:col with lines" > /tmp/blah.jpeg && open /tmp/blah.jpeg
//
func RunBackTest(db *models.DB, userId int) {

	// Build screener object
	screen := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 5.0},
			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.50},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
			{UserId: 1, Key: "allow-more-than-one-expire", Operator: "=", ValueString: "no"},
			{UserId: 1, Key: "allow-more-than-one-strike", Operator: "=", ValueString: "no"},
		},
	}

	// Set backtest
	btM := models.Backtest{
		UserId:          1,
		StartingBalance: 10000.00,
		EndingBalance:   10000.00,
		PositionSize:    "15-percent", // one-at-time, *-percent
		StartDate:       models.Date{helpers.ParseDateNoError("2017-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2020-12-31")},
		Midpoint:        false,
		TradeSelect:     "highest-credit",
		Benchmark:       "SPY",
		Screen:          screen,
	}

	// Setup a new backtesting
	bt := backtesting.New(db, btM.Benchmark)

	// Run blank backtest
	bt.DoBacktestDays(&btM)

	plotData := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	csvData := [][]string{{"Open Date", "Close Date", "Spread", "Open", "Close", "Lots", "% Away", "Margin", "Balance", "Return", "Benchmark", "Benchmark Balance", "Benchmark Return", "Status", "Note"}}
	table.SetHeader([]string{"Open Date", "Close Date", "Spread", "Open", "Close", "Lots", "% Away", "Margin", "Balance", "Return", "Benchmark", "Benchmark Balance", "Benchmark Return", "Status", "Note"})
	closedReturn := 0.00
	investedBenchmark := math.Floor(btM.StartingBalance / btM.BenchmarkStart)
	investedBenchmarkLeftOver := btM.StartingBalance - (investedBenchmark * btM.BenchmarkStart)

	for _, row := range btM.Positions {
		// Only record closed position
		if row.Status == "Closed" {
			closedReturn = row.ReturnFromStart
		}

		// Benchmark return
		bReturn := (((row.BenchmarkLast - btM.BenchmarkStart) / btM.BenchmarkStart) * 100)
		bAmountReturn := (investedBenchmark * row.BenchmarkLast) + investedBenchmarkLeftOver

		// Build data string
		d := []string{
			row.OpenDate.Format("01/02/2006"),
			row.CloseDate.Format("01/02/2006"),
			fmt.Sprintf("%s %s %.2f / %.2f", row.Legs[0].OptionUnderlying, row.Legs[0].OptionExpire.Format("01/02/2006"), row.Legs[0].OptionStrike, row.Legs[1].OptionStrike),
			fmt.Sprintf("$%.2f", row.OpenPrice),
			fmt.Sprintf("$%.2f", row.ClosePrice),
			fmt.Sprintf("%d", row.Lots),
			fmt.Sprintf("%.2f", row.PutPrecentAway) + "%",
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(row.Margin))),
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(row.Balance))),
			fmt.Sprintf("%.2f", row.ReturnFromStart) + "%",
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(row.BenchmarkLast))),
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(bAmountReturn))),
			fmt.Sprintf("%.2f", bReturn) + "%",
			row.Status,
			row.Note,
		}

		table.Append(d)
		csvData = append(csvData, d)
		plotData = append(plotData, []string{row.OpenDate.Format("01/02/2006"), humanize.BigCommaf(big.NewFloat(row.Balance))})
	}
	table.Render()

	// CAGR = (Ending value / Starting value)^(1 / # years) -1
	s1 := helpers.ParseDateNoError(btM.StartDate.Format("01/02/2006"))
	s2 := helpers.ParseDateNoError(btM.EndDate.Format("01/02/2006"))
	years := (s2.Sub(s1).Hours() / 24) / 365
	cagr := (math.Pow((btM.EndingBalance/btM.StartingBalance), (1/years)) - 1) * 100
	benchmarkCagr := (math.Pow((btM.BenchmarkEnd/btM.BenchmarkStart), (1/years)) - 1) * 100

	// Summary data
	tradeCount := len(btM.Positions)
	profit := (btM.EndingBalance - btM.StartingBalance)
	benchmarkPercent := (((btM.BenchmarkEnd - btM.BenchmarkStart) / btM.BenchmarkStart) * 100)

	// Show how long the backtest took.
	log.Printf("Backtest took %s", btM.TimeElapsed)
	log.Println("")
	log.Println("Summmary")
	log.Println("-------------")
	log.Printf("CAGR: %s%%", humanize.BigCommaf(big.NewFloat(cagr)))
	log.Printf("Return: %s%%", humanize.BigCommaf(big.NewFloat(closedReturn)))
	log.Printf("Profit: $%s", humanize.BigCommaf(big.NewFloat(profit)))
	log.Printf("Trade Count: %d", tradeCount)
	log.Println("")
	log.Println("Benchmark")
	log.Println("-------------")
	log.Printf("Start (%s): %s", btM.Benchmark, humanize.BigCommaf(big.NewFloat(btM.BenchmarkStart)))
	log.Printf("End (%s): %s", btM.Benchmark, humanize.BigCommaf(big.NewFloat(btM.BenchmarkEnd)))
	log.Printf("CAGR (%s): %s%%", btM.Benchmark, humanize.BigCommaf(big.NewFloat(benchmarkCagr)))
	log.Printf("Return (%s): %s%%", btM.Benchmark, humanize.BigCommaf(big.NewFloat(benchmarkPercent)))
	log.Println("")

	// ------------------ Export CSV ----------- //

	file, err := os.Create("/Users/spicer/Downloads/result.csv")

	if err != nil {
		spew.Dump(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range csvData {
		err := writer.Write(value)

		if err != nil {
			spew.Dump(err)
		}

	}

	// --------- Graph CSV -------- //

	file2, err := os.Create("/Users/spicer/Downloads/graph-balance.csv")

	if err != nil {
		spew.Dump(err)
	}

	defer file2.Close()

	writer = csv.NewWriter(file2)
	defer writer.Flush()

	for _, value := range plotData {
		err := writer.Write(value)

		if err != nil {
			spew.Dump(err)
		}

	}
}

/* End File */
