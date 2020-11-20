//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Delete a screener - 01
//
func TestDeleteScreener01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Post data
	screenerPost := models.Screener{
		UserId:   2,
		Name:     "Super Cool Test Screener - Deleted",
		Strategy: "put-credit-spread",
		Symbol:   "SPY",
		Items: []models.ScreenerItem{
			{Key: "spread-width", Operator: "=", ValueString: "", ValueNumber: 2.0},
			{Key: "max-days-to-expire", Operator: "=", ValueString: "", ValueNumber: 45},
			{Key: "short-strike-percent-away", Operator: "=", ValueString: "", ValueNumber: 5.0},
		},
	}

	db.Create(&screenerPost)

	// Make a mock request.
	req, _ := http.NewRequest("DELETE", "/api/v1/screeners/4", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.DELETE("/api/v1/screeners/:id", c.DeleteScreener)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Make sure the screener is deleted
	_, err := db.GetScreenerByIdAndUserId(uint(4), uint(2))

	// Validate result
	st.Expect(t, err.Error(), "[Models:GetScreenerByIdAndUserId] Record not found")
	st.Expect(t, w.Code, 204)
}

//
// UpdateScreener - 01
//
func TestUpdateScreener01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Post data
	screenerPost := models.Screener{
		UserId:   2,
		Name:     "Super Cool Test Screener",
		Strategy: "put-credit-spread",
		Symbol:   "SPY",
		Items: []models.ScreenerItem{
			{Key: "spread-width", Operator: "=", ValueString: "", ValueNumber: 2.0},
			{Key: "max-days-to-expire", Operator: "=", ValueString: "", ValueNumber: 45},
			{Key: "short-strike-percent-away", Operator: "=", ValueString: "", ValueNumber: 5.0},
		},
	}

	db.Create(&screenerPost)

	// Add update.
	screenerPost.Name = "Super Cool Test Screener - Update"

	screenerPost.Items = []models.ScreenerItem{
		{Key: "spread-width", Operator: "=", ValueString: "", ValueNumber: 2.0},
		{Key: "max-days-to-expire", Operator: "=", ValueString: "", ValueNumber: 45},
		{Key: "min-days-to-expire", Operator: "=", ValueString: "", ValueNumber: 5.0},
	}

	// Convert to JSON
	postStr, _ := json.Marshal(screenerPost)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/screeners/4", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.PUT("/api/v1/screeners/:id", c.UpdateScreener)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate result
	st.Expect(t, w.Code, 204)
}

//
// CreateScreener - 01
//
func TestCreateScreener01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Post data
	screenerPost := models.Screener{
		Name:     "Super Cool Test Screener",
		Strategy: "put-credit-spread",
		Symbol:   "SPY",
		Items: []models.ScreenerItem{
			{Key: "spread-width", Operator: "=", ValueString: "", ValueNumber: 2.0},
			{Key: "max-days-to-expire", Operator: "=", ValueString: "", ValueNumber: 45},
			{Key: "short-strike-percent-away", Operator: "=", ValueString: "", ValueNumber: 5.0},
		},
	}

	postStr, _ := json.Marshal(screenerPost)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/screeners", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/screeners", c.CreateScreener)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate result
	st.Expect(t, w.Body.String(), `{"id":4,"user_id":2,"name":"Super Cool Test Screener","strategy":"put-credit-spread","symbol":"SPY","items":[{"id":16,"screener_id":4,"user_id":2,"key":"spread-width","operator":"=","value_string":"","value_number":2},{"id":17,"screener_id":4,"user_id":2,"key":"max-days-to-expire","operator":"=","value_string":"","value_number":45},{"id":18,"screener_id":4,"user_id":2,"key":"short-strike-percent-away","operator":"=","value_string":"","value_number":5}]}`)
}

//
// GetScreenerResults01
//
func TestGetScreenerResults01(t *testing.T) {

	// // Start the db connection.
	// db, _ := models.NewDB()

	// // Create controller
	// c := &Controller{DB: db}

	// // Make a mock request.
	// req, _ := http.NewRequest("GET", "/api/v1/screeners/1/results", nil)
	// req.Header.Set("Accept", "application/json")

	// // Setup GIN Router
	// gin.SetMode("release")
	// gin.DisableConsoleColor()
	// r := gin.New()

	// r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	// r.GET("/api/v1/screeners/:id/results", c.GetScreenerResults)

	// // Setup writer.
	// w := httptest.NewRecorder()
	// r.ServeHTTP(w, req)

	// fmt.Println(w.Body.String())

}

//
// GetScreenerResultsFromFilters01
//
func TestGetScreenerResultsFromFilters01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Post data
	screenerPost := models.Screener{
		Name:     "One Time",
		Strategy: "put-credit-spread",
		Symbol:   "SPY",
		Items: []models.ScreenerItem{
			{Key: "min-credit", Operator: "=", ValueString: "", ValueNumber: 0.18},
			{Key: "spread-width", Operator: "=", ValueString: "", ValueNumber: 2.0},
			{Key: "max-days-to-expire", Operator: "=", ValueString: "", ValueNumber: 45},
			{Key: "short-strike-percent-away", Operator: "=", ValueString: "", ValueNumber: 4},
		},
	}

	postStr, _ := json.Marshal(screenerPost)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/screeners/results", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/screeners/results", c.GetScreenerResultsFromFilters)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Println(w.Body.String())

}

//
// TestGetScreeners01
//
func TestGetScreeners01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/screeners", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/screeners", c.GetScreeners)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := []models.Screener{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, len(result), 2)
	st.Expect(t, w.Body.String(), `[{"id":1,"user_id":1,"name":"Spicer's Dope Strategy","strategy":"put-credit-spread","symbol":"SPY","items":[{"id":1,"screener_id":1,"user_id":1,"key":"short-strike-percent-away","operator":"=","value_string":"","value_number":2.5},{"id":2,"screener_id":1,"user_id":1,"key":"spread-width","operator":"=","value_string":"","value_number":2},{"id":3,"screener_id":1,"user_id":1,"key":"min-credit","operator":"=","value_string":"","value_number":0.18},{"id":4,"screener_id":1,"user_id":1,"key":"max-days-to-expire","operator":"=","value_string":"","value_number":45},{"id":5,"screener_id":1,"user_id":1,"key":"min-days-to-expire","operator":"=","value_string":"","value_number":0}]},{"id":3,"user_id":1,"name":"Spicer's 2nd Dope Strategy","strategy":"iron-condor","symbol":"SPY","items":[{"id":11,"screener_id":3,"user_id":1,"key":"short-strike-percent-away","operator":"=","value_string":"","value_number":2.5},{"id":12,"screener_id":3,"user_id":1,"key":"spread-width","operator":"=","value_string":"","value_number":2},{"id":13,"screener_id":3,"user_id":1,"key":"min-credit","operator":"=","value_string":"","value_number":0.18},{"id":14,"screener_id":3,"user_id":1,"key":"max-days-to-expire","operator":"=","value_string":"","value_number":45},{"id":15,"screener_id":3,"user_id":1,"key":"min-days-to-expire","operator":"=","value_string":"","value_number":0}]}]`)
}

//
// TestGetScreener01 - Get a screener by id. - 01
//
func TestGetScreener01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/screeners/1", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/screeners/:id", c.GetScreener)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Screener{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Strategy, "put-credit-spread")
	st.Expect(t, len(result.Items), 5)
	st.Expect(t, w.Body.String(), `{"id":1,"user_id":1,"name":"Spicer's Dope Strategy","strategy":"put-credit-spread","symbol":"SPY","items":[{"id":1,"screener_id":1,"user_id":1,"key":"short-strike-percent-away","operator":"=","value_string":"","value_number":2.5},{"id":2,"screener_id":1,"user_id":1,"key":"spread-width","operator":"=","value_string":"","value_number":2},{"id":3,"screener_id":1,"user_id":1,"key":"min-credit","operator":"=","value_string":"","value_number":0.18},{"id":4,"screener_id":1,"user_id":1,"key":"max-days-to-expire","operator":"=","value_string":"","value_number":45},{"id":5,"screener_id":1,"user_id":1,"key":"min-days-to-expire","operator":"=","value_string":"","value_number":0}]}`)
}

//
// TestGetScreener02 - Get a screener by id. - 02 - Not Found
//
func TestGetScreener02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/screeners/4", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/screeners/:id", c.GetScreener)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"No Record Found."}`)
}

//
// TestGetScreener03 - Get a screener by id. - 03 - No Perms
//
func TestGetScreener03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/screeners/2", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/screeners/:id", c.GetScreener)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"No Record Found."}`)
}

/* End File */
