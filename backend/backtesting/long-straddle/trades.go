//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package longstraddle

import (
	"fmt"
	"os"
	"time"

	"app.options.cafe/backtesting/helpers"
	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"app.options.cafe/screener"
	"github.com/davecgh/go-spew/spew"
)

//
// Trades will place our trades
//
func Trades(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem, underlyingLast float64) {
	// First look if there are any trades we need to close.
	Close(today, underlyingLast, backtest, options)

	// Figure out qtys. QTY field in symbol is not stored in DB. Just used for this purpose
	results[0].Legs[0].Qty = 1
	results[0].Legs[1].Qty = 1

	// Open the trade
	helpers.OpenTrade(today, backtest, results[0])

	spew.Dump(backtest)

	os.Exit(1)

	// for _, row := range results {

	// 	upperPrice := row.UnderlyingLast + row.Ask
	// 	breakEvenJump := ((upperPrice - row.UnderlyingLast) / row.UnderlyingLast) * 100

	// 	fmt.Println(row.Ask, " : ", breakEvenJump)
	// }

	fmt.Println(today, " ", "Trades")

	return
}
