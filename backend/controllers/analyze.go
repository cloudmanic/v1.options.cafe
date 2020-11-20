//
// Date: 2018-11-16
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-16
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"app.options.cafe/library/analyze"
	"github.com/gin-gonic/gin"
)

//
// AnalyzeOptionsProfitLossByUnderlyingPrice
// Return data for analyzing a trade by underlying price at expire
//
func (t *Controller) AnalyzeOptionsProfitLossByUnderlyingPrice(c *gin.Context) {

	var trade analyze.Trade

	// Json to Object
	if err := c.ShouldBindJSON(&trade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Loop through and add the symbol to the object
	for key, row := range trade.Legs {

		sym, err := t.DB.CreateNewOptionSymbol(row.SymbolStr)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No symbol found."})
			return
		}

		trade.Legs[key].Symbol = sym
	}

	// Get the Profit and Loss By Underlying Price
	results := analyze.OptionsProfitLossByUnderlyingPrice(trade)

	// Return happy JSON
	c.JSON(200, results)
}

/* End File */
