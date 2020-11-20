//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-18
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"app.options.cafe/brokers/tradier"
	"app.options.cafe/brokers/types"
	"app.options.cafe/library/cache"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/services"
	"github.com/gin-gonic/gin"
)

//
// Get quotes
//
func (t *Controller) GetQuotes(c *gin.Context) {

	quotes := []types.Quote{}

	smbs := strings.Split(c.Query("symbols"), ",")

	// First see if we have all these quotes in cache?
	allCachedFound := true

	for _, row := range smbs {

		cachedQuote := types.Quote{}

		found, _ := cache.Get("oc-paper-trade-quote-"+strings.ToLower(row), &cachedQuote)

		// Return happy JSON
		if !found {
			allCachedFound = false
		} else {
			quotes = append(quotes, cachedQuote)
		}

	}

	// All the quotes we requested were not found in cache. Make API to Tradier to update.
	if !allCachedFound {

		// Setup the broker
		broker := tradier.Api{
			DB:     masterDB,
			ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"),
		}

		// Get quotes
		quotes, err := broker.GetQuotes(smbs)

		if err != nil {
			services.Info(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get quotes."})
			return
		}

		// Store the quotes in cache
		for _, row := range quotes {
			cache.SetExpire("oc-paper-trade-quote-"+strings.ToLower(row.Symbol), (time.Minute * 15), row)
		}

	}

	// Loop through and change the values because this is paper trading and all.
	for key, row := range quotes {
		r := rand.Float64()
		quotes[key].Last = helpers.Round(row.Last+r, 2)
	}

	// Return happy JSON
	c.JSON(200, quotes)
}

/* End File */
