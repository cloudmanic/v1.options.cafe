//
// Date: 11/09/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/gin-gonic/gin"
)

//
// Return symbols in our database.
//
func (t *Controller) GetSymbols(c *gin.Context) {

	// Search for symbol
	if c.Query("search") != "" {
		t.DoSymbolSearch(c)
	}
}

//
// Do Symbol Search
//
func (t *Controller) DoSymbolSearch(c *gin.Context) {

	// Get the query.
	search := c.Query("search")

	// Run DB query
	symbols, err := t.DB.SearchSymbols(search)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Return happy JSON
	c.JSON(200, symbols)
}

/* End File */
