//
// Date: 11/09/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/araddon/dateparse"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
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
	symbols, err := t.DB.SearchSymbols(search, "Equity")

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": httpGenericErrMsg})
		return
	}

	// Return happy JSON
	c.JSON(200, symbols)
}

//
// Add a symbol to active symbols by user.
//
func (t *Controller) AddActiveSymbol(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	symbol := gjson.Get(string(body), "symbol").String()

	// Validate name
	if len(symbol) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol field can not be empty."})
		return
	}

	// Store the symbol
	act, err := t.DB.CreateActiveSymbol(userId, symbol)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": httpGenericErrMsg})
		return
	}

	// Return happy JSON
	c.JSON(200, act)
}

//
// Post in Root Symbol, Expire, Strike, Type. Return a symbol.
//
func (t *Controller) GetOptionSymbolFromParts(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Parse the post.
	symbol := gjson.Get(string(body), "symbol").String()
	expire := gjson.Get(string(body), "expire").String()
	strike := gjson.Get(string(body), "strike").Float()
	optionType := gjson.Get(string(body), "type").String()

	// Set expire date
	expireDate, err := dateparse.ParseAny(expire)

	if t.RespondError(c, err, "Invalid expire date.") {
		return
	}

	// Query database
	symb, err := t.DB.GetOptionByParts(symbol, optionType, expireDate, strike)

	if t.RespondError(c, err, "Unable to find symbol.") {
		return
	}

	// Now add this symbol to our active symbol list as we most likely want quotes via websockets after this.
	_, err2 := t.DB.CreateActiveSymbol(userId, symb.ShortName)

	if err2 != nil {
		services.BetterError(err2)
		c.JSON(http.StatusBadRequest, gin.H{"error": httpGenericErrMsg})
		return
	}

	// Return happy JSON
	c.JSON(200, symb)
}

/* End File */
