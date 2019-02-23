//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

const cacheDirBase = "backtesting-options-chains"

// Base struct
type Base struct {
	DB            models.Datastore
	StrategyFuncs map[string]func(today time.Time, backtest models.Backtest, underlyingLast float64, chains map[time.Time]types.OptionsChain) error
}

//
// New Backtest
//
func New(db models.Datastore) Base {

	// New backtest instance
	t := Base{
		DB: db,
	}

	// Build backtest functions.
	t.StrategyFuncs = map[string]func(today time.Time, backtest models.Backtest, underlyingLast float64, chains map[time.Time]types.OptionsChain) error{
		"blank":             t.DoBlank,
		"put-credit-spread": t.DoPutCreditSpread,
	}

	// Set the cache dir.
	cacheDir := os.Getenv("CACHE_DIR") + "/" + cacheDirBase

	// Make a directory to download.
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}

	return t
}

//
// DoBacktestDays - Loop through every day in the backtest and pass
// an options chain into a function for a prticular backtest type.
//
func (t *Base) DoBacktestDays(backtest models.Backtest) error {

	// Get dates by symbol
	dates, err := eod.GetTradeDatesBySymbols(backtest.Screen.Symbol)

	if err != nil {
		return err
	}

	// Loop through the dates and run backtest.
	for _, row := range dates {

		// We skip dates before our start date
		if row.Before(backtest.StartDate) {
			continue
		}

		// We skip dates after our end date
		if row.After(backtest.EndDate) {
			continue
		}

		// Create broker object
		o := eod.Api{
			DB:  t.DB,
			Day: row,
		}

		// Log where we are in the backtest. TODO(spicer): Send this up websocket
		//services.Info("Backtesting " + backtest.Screen.Strategy + " " + backtest.Screen.Symbol + " on " + row.Format("2006-01-02"))

		// See if we have the chain in cache
		chains, err := getCachedChain(backtest.Screen.Symbol, row)

		// We did not have the chains in the file cache
		if err != nil {
			// Log no cache found.
			services.Info("Backtesting - chains not found in file cache - " + backtest.Screen.Symbol + " " + row.Format("2006-01-02"))

			// Get the expire dates for this option option chain.
			expDates, err2 := o.GetOptionsExpirationsBySymbol(backtest.Screen.Symbol)

			if err2 != nil {
				return err2
			}

			// Create map of expire dates and thier chain
			tmpChains := map[time.Time]types.OptionsChain{}

			// Loop through expire dates and look for possible trades
			for _, row2 := range expDates {

				// Get the options change by expire
				chain, err2 := o.GetOptionsChainByExpiration(backtest.Screen.Symbol, row2)

				if err2 != nil {
					return err2
				}

				// Add to chains map
				tmpChains[helpers.ParseDateNoError(row2)] = chain
			}

			// Store in file cache
			setCacheChain(backtest.Screen.Symbol, row, tmpChains)

			// Reset chains
			chains = tmpChains
		}

		// Get underlyingLast
		underlyingLast := 0.00
		for _, row2 := range chains {
			underlyingLast = row2.UnderlyingLast
		}

		// Run backtest strategy function for this backtest
		err = t.StrategyFuncs[backtest.Screen.Strategy](row, backtest, underlyingLast, chains)

		if err != nil {
			return err
		}
	}

	return nil
}

//
// DoBlank is mostly using for unit testing.
//
func (t *Base) DoBlank(today time.Time, backtest models.Backtest, underlyingLast float64, chains map[time.Time]types.OptionsChain) error {
	return nil
}

//
// setCacheChain store a chance in file cache
//
func setCacheChain(symbol string, today time.Time, chains map[time.Time]types.OptionsChain) {
	// Set the cache dir.
	cacheDir := os.Getenv("CACHE_DIR") + "/" + cacheDirBase

	// Store results in file cache.
	j, err := json.Marshal(chains)

	if err != nil {
		return
	}

	err = ioutil.WriteFile(cacheDir+"/"+symbol+"-"+today.Format("2006-01-02")+".json", j, 0644)

	if err != nil {
		return
	}
}

//
// getCachedChain - See if we have a cached chain stored on file.
//
func getCachedChain(symbol string, today time.Time) (map[time.Time]types.OptionsChain, error) {
	// Set the cache dir.
	cacheDir := os.Getenv("CACHE_DIR") + "/" + cacheDirBase
	cacheFile := cacheDir + "/" + symbol + "-" + today.Format("2006-01-02") + ".json"

	// See if we have the file.
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, errors.New("cache not found")
	}

	// Read contents of file.
	fileDat, err := ioutil.ReadFile(cacheFile)

	if err != nil {
		return nil, err
	}

	// JSON to struct
	var dat map[time.Time]types.OptionsChain

	if err := json.Unmarshal(fileDat, &dat); err != nil {
		return nil, err
	}

	// Retun happy
	return dat, nil
}

/* End File */
