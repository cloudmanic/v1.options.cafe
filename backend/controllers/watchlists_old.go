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
// Test UpdateWatchlist - 01
//
func TestUpdateWatchlist01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Body data
	var bodyStr = []byte(`{"name":"An Updated Super Cool Name"}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/watchlists/3", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.PUT("/api/v1/watchlists/:id", c.UpdateWatchlist)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate response
	st.Expect(t, w.Code, 204)

	// Now verify the watchlist was deleted.
	wList, err := db.GetWatchlistsById(3)

	// Validate
	st.Expect(t, err, nil)
	st.Expect(t, wList.Id, uint(3))
	st.Expect(t, wList.Name, "An Updated Super Cool Name")
}

//
// Test UpdateWatchlist - 02 (No access)
//
func TestUpdateWatchlist02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Body data
	var bodyStr = []byte(`{"name":"An Updated Super Cool Name"}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/watchlists/2", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.PUT("/api/v1/watchlists/:id", c.UpdateWatchlist)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate response
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"No Record Found."}`)

	// Now verify the watchlist was not changed.
	wList, err := db.GetWatchlistsById(2)

	// Validate
	st.Expect(t, err, nil)
	st.Expect(t, wList.Id, uint(2))
	st.Expect(t, wList.Name, "Watchlist #2 - User 1")
}

//
// Test UpdateWatchlist - 03 (No list found)
//
func TestUpdateWatchlist03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Body data
	var bodyStr = []byte(`{"name":"An Updated Super Cool Name"}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/watchlists/200", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.PUT("/api/v1/watchlists/:id", c.UpdateWatchlist)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate response
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"No Record Found."}`)
}

//
// Test DeleteWatchlist - 01
//
func TestDeleteWatchlist01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("DELETE", "/api/v1/watchlists/3", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.DELETE("/api/v1/watchlists/:id", c.DeleteWatchlist)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate response
	st.Expect(t, w.Code, 204)

	// Now verify the watchlist was deleted.
	wList, err := db.GetWatchlistsById(3)

	// Validate
	st.Expect(t, err.Error(), "[Models:GetWatchlistsById] Record not found")
	st.Expect(t, wList.Id, uint(0))

	// Verify the symbols were also deleted
	sList := db.WatchlistSymbolGetByWatchlistId(3)

	// Validate
	st.Expect(t, len(sList), 0)
}

//
// Test CreateWatchlist - 01
//
func TestCreateWatchlist01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{"name":"Super Cool Test Watchlist"}`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/watchlists", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/watchlists", c.CreateWatchlist)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Watchlist{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, result.Id, uint(5))
	st.Expect(t, result.Name, "Super Cool Test Watchlist")
	st.Expect(t, result.UserId, uint(2))
	st.Expect(t, len(result.Symbols), 0)
}

//
// Test CreateWatchlist - 02 (empty name)
//
func TestCreateWatchlist02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{"name":""}`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/watchlists", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/watchlists", c.CreateWatchlist)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Body.String(), `{"error":"Name field can not be empty."}`)
}

//
// Test CreateWatchlist - 03 (duplicate name)
//
func TestCreateWatchlist03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{"name":"Watchlist #2 - User 2"}`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/watchlists", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/watchlists", c.CreateWatchlist)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Body.String(), `{"error":"A watchlist already exists by this name."}`)
}

//
// Test - WatchlistAddSymbol - 01 (Success)
//
func TestWatchlistAddSymbol01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

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

	// ----------- Verify the order in the DB is correct --------------- //

	// List of watchlists.
	list := []models.WatchlistSymbol{}

	// Loop through the watch list and move the order of each symbol down by one in order.
	db.Order("`order` asc").Where("watchlist_id = ?", result.WatchlistId).Find(&list)

	// Test results.
	st.Expect(t, len(list), 4)
	st.Expect(t, list[0].Id, uint(13))
	st.Expect(t, list[0].Order, uint(0))
}

//
// Test - WatchlistAddSymbol - 02 (Another user's watchlist id)
//
func TestWatchlistAddSymbol02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

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
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

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
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

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

//
// Test WatchlistReorder - 01
//
func TestWatchlistReorder01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Body data
	var bodyStr = []byte(`{ "ids": [ 8, 9, 7 ] }`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/watchlists/3/reorder", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.PUT("/api/v1/watchlists/:id/reorder", c.WatchlistReorder)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 204)

	// Check the database to make sure the ids are in the correct order.
	wList, err := db.GetWatchlistsById(3)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, wList.Id, uint(3))
	st.Expect(t, wList.Symbols[0].Id, uint(8))
	st.Expect(t, wList.Symbols[1].Id, uint(9))
	st.Expect(t, wList.Symbols[2].Id, uint(7))
}

//
// Test WatchlistReorder - 02 (No Access)
//
func TestWatchlistReorder02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Body data
	var bodyStr = []byte(`{ "ids": [ 8, 9, 7 ] }`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/watchlists/2/reorder", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.PUT("/api/v1/watchlists/:id/reorder", c.WatchlistReorder)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 401)
	st.Expect(t, w.Body.String(), `{"error":"No access to this watchlist resource."}`)
}

//
// Test WatchlistDeleteSymbol - 01
//
func TestWatchlistDeleteSymbol01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("DELETE", "/api/v1/watchlists/3/symbol/8", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.DELETE("/api/v1/watchlists/:id/symbol/:symb", c.WatchlistDeleteSymbol)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 204)

	// Check the database to make sure the ids are in the correct order.
	wList, err := db.GetWatchlistsById(3)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, wList.Id, uint(3))
	st.Expect(t, len(wList.Symbols), 2)
	st.Expect(t, wList.Symbols[0].Id, uint(7))
	st.Expect(t, wList.Symbols[1].Id, uint(9))
}

//
// Test WatchlistDeleteSymbol - 02 (No Access)
//
func TestWatchlistDeleteSymbol02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("DELETE", "/api/v1/watchlists/2/symbol/8", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.DELETE("/api/v1/watchlists/:id/symbol/:symb", c.WatchlistDeleteSymbol)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 401)
	st.Expect(t, w.Body.String(), `{"error":"No access to this watchlist resource."}`)
}

/* End File */
