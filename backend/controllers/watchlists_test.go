//
// Date: 3/5/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"github.com/tidwall/gjson"
)

//
// Test - WatchlistAddSymbol - 01 (Success)
//
func TestWatchlistAddSymbol01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{ "symbol_id": 7 }`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/watchlists/3/symbol", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/watchlists/:id/symbol", c.WatchlistAddSymbol)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.WatchlistSymbol{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, result.UserId, uint(0)) // We remove this in json
	st.Expect(t, result.WatchlistId, uint(3))
	st.Expect(t, result.SymbolId, uint(7))
	st.Expect(t, result.Order, uint(0))
	st.Expect(t, result.Symbol.ShortName, "CAT")
}

//
// Test - WatchlistAddSymbol - 02 (Another user's watchlist id)
//
func TestWatchlistAddSymbol02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{ "symbol_id": 7 }`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/watchlists/2/symbol", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/watchlists/:id/symbol", c.WatchlistAddSymbol)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 401)
	st.Expect(t, gjson.Get(w.Body.String(), "error").String(), "No access to this watchlist resource.")
}

//
// Test - WatchlistAddSymbol - 03 (Invalid symbol id)
//
func TestWatchlistAddSymbol03(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{ "symbol_id": 7999 }`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/watchlists/3/symbol", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/watchlists/:id/symbol", c.WatchlistAddSymbol)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.symbol_id").String(), "Unknown symbol_id.")
}

//
// Test - WatchlistAddSymbol - 04 (A symbol that is already added to the watchlist)
//
func TestWatchlistAddSymbol04(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{ "symbol_id": 1 }`) // Already added in the testing

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/watchlists/3/symbol", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/watchlists/:id/symbol", c.WatchlistAddSymbol)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "error").String(), "Symbol already part of this watchlist.")
}

/* End File */
