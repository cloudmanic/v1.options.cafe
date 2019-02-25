//
// Date: 2018-10-30
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-07
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package eod

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

//
// GetOptionsExpirationsBySymbol - Get options expiration by Symbol
//
func (t *Api) GetOptionsExpirationsBySymbol(symbol string) ([]string, error) {

	expires := []string{}
	symb := strings.ToUpper(symbol)
	tmpExpires := make(map[string]bool)
	cacheKey := "oc-brokers-eod-expirations-" + strings.ToUpper(symbol) + "-" + t.Day.Format("2006-01-02")

	// See if we have this result in the cache.
	var cacheddates []string
	found, _ := cache.Get(cacheKey, &cacheddates)

	// Return happy JSON
	if found {
		return cacheddates, nil
	}

	// Get a list of all options
	options, _, err := t.GetOptionsBySymbol(symb)

	if err != nil {
		return expires, err
	}

	// Loop through and build chain
	for _, row := range options {
		tmpExpires[row.ExpirationDate.Format("2006-01-02")] = true
	}

	// Loop through and create an array of string
	for key := range tmpExpires {
		expires = append(expires, key)
	}

	// Sort the results date in asc order.
	sort.Strings(expires)

	// Store dates in cache. Does not expire.
	cache.Set(cacheKey, expires)

	// Return happy
	return expires, nil
}

//
// GetOptionsChainByExpiration - Get an options chain by expiration.
//
func (t *Api) GetOptionsChainByExpiration(symbol string, expireStr string) (types.OptionsChain, error) {

	symb := strings.ToUpper(symbol)
	expireDate := types.Date{helpers.ParseDateNoError(expireStr).UTC()}

	// Get a list of all options
	options, underlyingLast, err := t.GetOptionsBySymbol(symb)

	// New chain to return
	chain := types.OptionsChain{
		Underlying:     symb,
		UnderlyingLast: underlyingLast,
		ExpirationDate: expireDate,
		Puts:           []types.OptionsChainItem{},
		Calls:          []types.OptionsChainItem{},
	}

	if err != nil {
		return chain, err
	}

	// Loop through and build chain
	for _, row := range options {

		// We only want the expire date we passed in.
		if row.ExpirationDate != expireDate {
			continue
		}

		// Append Item
		if row.OptionType == "Call" {
			chain.Calls = append(chain.Calls, row)
		} else if row.OptionType == "Put" {
			chain.Puts = append(chain.Puts, row)
		}

	}

	// Sort Strikes. - Calls
	sort.Slice(chain.Calls, func(i, j int) bool {
		return chain.Calls[i].Strike < chain.Calls[j].Strike
	})

	// Sort Strikes. - Puts
	sort.Slice(chain.Puts, func(i, j int) bool {
		return chain.Puts[i].Strike < chain.Puts[j].Strike
	})

	// Store in sql cache.
	setSqlLiteChain(symb, t.Day, chain)

	// Return Chain
	return chain, nil
}

//
// Store this chain in a file cache so we can get it faster in the future.
//
func setSqlLiteChain(symbol string, today time.Time, chain types.OptionsChain) {

	// Set the cache dir and sqlfile
	cacheDir := os.Getenv("CACHE_DIR") + "/" + cacheDirBase + "/chain/" + symbol
	dbFile := cacheDir + "/" + today.Format("2006-01-02") + ".sqlite"

	fmt.Println(dbFile)

	// Make a directory to create sqlite db in.
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}

	// Connect to sqlite db.
	db, err := gorm.Open("sqlite3", dbFile)

	if err != nil {
		services.Fatal(errors.New("setSqlLiteChain: failed to connect sqlite database - " + dbFile))
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&types.OptionsChain{})
	db.AutoMigrate(&types.OptionsChainItem{})

	// Create
	db.Create(&chain)

}

/* End File */
