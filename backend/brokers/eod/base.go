//
// Date: 2018-10-29
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-29
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// About: This is a broker that brings in data from eod archived data. Useful for back testing, and unit testing.
//

package eod

import (
	"go/build"
	"strings"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/store/object"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	env "github.com/jpfuentes2/go-env"
)

type Api struct {
	DB models.Datastore
}

//
// Start up the controller.
//
func init() {
	// Helpful for testing
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/.env")
}

//
// Return a list of trade dates by symbol
//
func GetTradeDatesBySymbols(symbol string) ([]time.Time, error) {

	dates := []time.Time{}
	dirPath := "options-eod/" + strings.ToUpper(symbol) + "/"
	cacheKey := "oc-brokers-eod-trade-dates-" + strings.ToUpper(symbol)

	// See if we have this result in the cache.
	cacheddates := []time.Time{}
	found, _ := cache.Get(cacheKey, &cacheddates)

	// Return happy JSON
	if found {
		return cacheddates, nil
	}

	// Make call to S3 store and get the dates for different EODs
	files, err := object.ListObjects(dirPath)

	if err != nil {
		return dates, err
	}

	// Convert strings to dates
	for _, row := range files {

		fileName := strings.Replace(row.Key, dirPath, "", -1)
		fileName = strings.Replace(fileName, ".csv.zip", "", -1)
		dates = append(dates, helpers.ParseDateNoError(fileName).UTC())

	}

	// Store dates in cache - 3 hours
	cache.SetExpire(cacheKey, (time.Minute * 60 * 3), dates)

	// Return happy
	return dates, nil
}

//
// Return a list of keys to download by symbol.
//
func GetTradeDateKeysBySymbol(symbol string) ([]string, error) {

	keys := []string{}
	dirPath := "options-eod/" + strings.ToUpper(symbol) + "/"

	// Get dates from S3 store
	dates, err := GetTradeDatesBySymbols("spy")

	if err != nil {
		return keys, err
	}

	// Loop through and create keys
	for _, row := range dates {
		keys = append(keys, dirPath+row.Format("2006-01-02")+".csv.zip")
	}

	// Return happy
	return keys, nil
}

/* End File */
