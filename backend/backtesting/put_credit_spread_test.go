//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	humanize "github.com/dustin/go-humanize"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nbio/st"
	"github.com/olekukonko/tablewriter"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

//
// TestRunPutCreditSpread01 - Run a put credit spread backtest.
//
func TestDoPutCreditSpread01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing")
	defer models.TestingTearDown(db, dbName)

	// Setup a new backtesting
	bt := New(db)

	// Build screener object
	screen := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 4.0},
			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.20},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// Set backtest
	btM := models.Backtest{
		UserId:          1,
		StartingBalance: 2000.00,
		EndingBalance:   2000.00,
		PositionSize:    "10-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2018-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2019-01-01")},
		Midpoint:        true,
		TradeSelect:     "highest-credit",
		Screen:          screen,
	}

	// Run blank backtest
	err := bt.DoBacktestDays(&btM)
	st.Expect(t, err, nil)

	//spew.Dump(btM)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Open Date", "Close Date", "Spread", "Open", "Close", "Lots", "% Away", "Margin", "Balance", "Status", "Note"})

	for _, row := range btM.Positions {
		table.Append([]string{
			row.OpenDate.Format("01/02/2006"),
			row.CloseDate.Format("01/02/2006"),
			fmt.Sprintf("%s %s %.2f / %.2f", row.Legs[0].OptionUnderlying, row.Legs[0].OptionExpire.Format("01/02/2006"), row.Legs[0].OptionStrike, row.Legs[1].OptionStrike),
			fmt.Sprintf("$%.2f", row.OpenPrice),
			fmt.Sprintf("$%.2f", row.ClosePrice),
			fmt.Sprintf("%d", row.Lots),
			fmt.Sprintf("%.2f", row.PutPrecentAway) + "%",
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(row.Margin))),
			fmt.Sprintf("$%s", humanize.BigCommaf(big.NewFloat(row.Balance))),
			row.Status,
			row.Note,
		})
	}
	table.Render()

	// Show how long the backtest took.
	log.Printf("Backtest took %s", btM.TimeElapsed)

}

/* End File */
