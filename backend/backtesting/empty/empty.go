//
// Date: 2022-05-02
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package empty

import (
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"app.options.cafe/screener"

	screenerCache "app.options.cafe/screener/cache"
)

//
// EmptyTrades is used for testing.
//
func Trades(db models.Datastore, today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem, underlyingLast float64, benchmarkQuotes []types.HistoryQuote) {

	return
}

//
// Results - Used for testing.
//
func Results(db models.Datastore, today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error) {
	// Results that we return.
	results := []screener.Result{}

	// Return happy with results.
	return results, nil
}

/* End File */
