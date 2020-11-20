//
// Date: 11/09/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/araddon/dateparse"
	"app.options.cafe/brokers/tradier"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//
// Return symbol in our database.
//
func (t *Controller) GetSymbol(c *gin.Context) {

	// Run DB query
	symbol := models.Symbol{}

	// Another test to see if search works
	err := t.DB.Query(&symbol, models.QueryParam{
		Wheres: []models.KeyValue{{Key: "short_name", Value: c.Param("symb")}},
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No symbol found."})
		return
	}

	// Return happy JSON
	c.JSON(200, symbol)
}

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
		services.Info(err)
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
		services.Info(err)
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

	// To save some processing lets see if we already have this symbol
	sym, err := t.DB.GetOptionByParts(symbol, optionType, expireDate, strike)

	if err == nil {

		// Now add this symbol to our active symbol list as we most likely want quotes via websockets after this.
		_, err2 := t.DB.CreateActiveSymbol(userId, sym.ShortName)

		if err2 != nil {
			services.Info(err2)
			c.JSON(http.StatusBadRequest, gin.H{"error": httpGenericErrMsg})
			return
		}

		// Return happy JSON
		c.JSON(200, sym)

		return
	}

	// ------------ First we load the entire option chain from Tradier : TODO: someday make this via the users' account -------------- //
	// ------------ We do this to make sure our symbols table is up to date with the latest for this symbol -------------------------- //

	// // Get access token
	// apiKey, err := t.GetTradierAccessToken(c)

	// if err != nil {
	// 	t.RespondError(c, err, httpGenericErrMsg)
	// 	return
	// }

	// Setup the broker
	broker := tradier.Api{
		DB:     t.DB,
		ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"),
	}

	// Get chain from tradier.
	chain, err := broker.GetOptionsChainByExpiration(symbol, expire)

	if err != nil {
		services.Info(err)
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Make suser the option we care about is in our Symbol DB.
	if optionType == "Call" {
		for _, row := range chain.Calls {
			if row.Strike == strike {
				t.DB.CreateNewOptionSymbol(row.Symbol)
			}
		}
	} else {
		for _, row := range chain.Puts {
			if row.Strike == strike {
				t.DB.CreateNewOptionSymbol(row.Symbol)
			}
		}
	}

	// Since we made the chain call lets tore the symbols
	go t.DB.LoadSymbolsByOptionsChain(chain)

	// ----------- Now Search for the option in the symbols table ------------ //

	// Query database
	symb, err := t.DB.GetOptionByParts(symbol, optionType, expireDate, strike)

	if t.RespondError(c, err, "Unable to find symbol.") {
		return
	}

	// Now add this symbol to our active symbol list as we most likely want quotes via websockets after this.
	_, err2 := t.DB.CreateActiveSymbol(userId, symb.ShortName)

	if err2 != nil {
		services.Info(err2)
		c.JSON(http.StatusBadRequest, gin.H{"error": httpGenericErrMsg})
		return
	}

	// Return happy JSON
	c.JSON(200, symb)
}

/* End File */
