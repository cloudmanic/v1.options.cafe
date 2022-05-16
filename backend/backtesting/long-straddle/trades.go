//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package longstraddle

import (
	"fmt"
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"app.options.cafe/screener"
)

//
// Trades will place our trades
//
func Trades(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem) {

	for _, row := range results {

		upperPrice := row.UnderlyingLast + row.Ask
		breakEvenJump := ((upperPrice - row.UnderlyingLast) / row.UnderlyingLast) * 100

		fmt.Println(row.Ask, " : ", breakEvenJump)
	}

	fmt.Println(today, " ", "Trades")

	return
}
