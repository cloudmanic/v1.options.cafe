//
// Date: 2022-05-02
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

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
func (t *Base) EmptyTrades(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem) {

	return
}

//
// EmptyResults - Used for testing.
//
func (t *Base) EmptyResults(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error) {
	// Results that we return.
	results := []screener.Result{}

	// Return happy with results.
	return results, nil
}

/* End File */
