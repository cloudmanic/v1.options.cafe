//
// Date: 4/12/2022
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestGetBacktests01 returns a list of backtests for a user.
//
func TestGetBacktests01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Load sample data into the database
	screen1 := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 4.0},
			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.20},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	screen2 := models.Screener{
		UserId:   1,
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{UserId: 1, Key: "short-strike-percent-away", Operator: ">", ValueNumber: 4.0},
			{UserId: 1, Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{UserId: 1, Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{UserId: 1, Key: "open-credit", Operator: "<", ValueNumber: 0.20},
			{UserId: 1, Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{UserId: 1, Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// Set backtest
	btM1 := models.Backtest{
		UserId:          1,
		StartingBalance: 2000.00,
		EndingBalance:   2000.00,
		PositionSize:    "10-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2017-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2021-12-31")},
		Midpoint:        true,
		TradeSelect:     "shortest-percent-away",
		Benchmark:       "SPY",
		Screen:          screen1,
	}

	// Set backtest
	btM2 := models.Backtest{
		UserId:          1,
		StartingBalance: 2000.00,
		EndingBalance:   2000.00,
		PositionSize:    "20-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2017-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2021-12-31")},
		Midpoint:        true,
		TradeSelect:     "shortest-percent-away",
		Benchmark:       "XLF",
		Screen:          screen2,
	}

	db.Save(&btM1)
	db.Save(&btM2)

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/backtests", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/backtests", c.GetBacktests)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	st.Expect(t, w.Body.String(), `[{"id":1,"user_id":1,"start_date":"2017-01-01","end_date":"2021-12-31","ending_balance":2000,"starting_balance":2000,"cagr":0,"return":0,"profit":0,"trade_count":0,"trade_select":"shortest-percent-away","midpoint":true,"position_size":"10-percent","time_elapsed":0,"benchmark":"SPY","benchmark_start":0,"benchmark_end":0,"benchmark_cagr":0,"benchmark_percent":0,"screen":{"id":1,"user_id":1,"backtest_id":1,"name":"","strategy":"put-credit-spread","symbol":"SPY","items":[{"id":1,"screener_id":1,"user_id":1,"key":"short-strike-percent-away","operator":"\u003e","value_string":"","value_number":4},{"id":2,"screener_id":1,"user_id":1,"key":"spread-width","operator":"=","value_string":"","value_number":2},{"id":3,"screener_id":1,"user_id":1,"key":"open-credit","operator":"\u003e","value_string":"","value_number":0.18},{"id":4,"screener_id":1,"user_id":1,"key":"open-credit","operator":"\u003c","value_string":"","value_number":0.2},{"id":5,"screener_id":1,"user_id":1,"key":"days-to-expire","operator":"\u003c","value_string":"","value_number":46},{"id":6,"screener_id":1,"user_id":1,"key":"days-to-expire","operator":"\u003e","value_string":"","value_number":0}]},"trade_groups":[]},{"id":2,"user_id":1,"start_date":"2017-01-01","end_date":"2021-12-31","ending_balance":2000,"starting_balance":2000,"cagr":0,"return":0,"profit":0,"trade_count":0,"trade_select":"shortest-percent-away","midpoint":true,"position_size":"20-percent","time_elapsed":0,"benchmark":"XLF","benchmark_start":0,"benchmark_end":0,"benchmark_cagr":0,"benchmark_percent":0,"screen":{"id":2,"user_id":1,"backtest_id":2,"name":"","strategy":"put-credit-spread","symbol":"SPY","items":[{"id":7,"screener_id":2,"user_id":1,"key":"short-strike-percent-away","operator":"\u003e","value_string":"","value_number":4},{"id":8,"screener_id":2,"user_id":1,"key":"spread-width","operator":"=","value_string":"","value_number":2},{"id":9,"screener_id":2,"user_id":1,"key":"open-credit","operator":"\u003e","value_string":"","value_number":0.18},{"id":10,"screener_id":2,"user_id":1,"key":"open-credit","operator":"\u003c","value_string":"","value_number":0.2},{"id":11,"screener_id":2,"user_id":1,"key":"days-to-expire","operator":"\u003c","value_string":"","value_number":46},{"id":12,"screener_id":2,"user_id":1,"key":"days-to-expire","operator":"\u003e","value_string":"","value_number":0}]},"trade_groups":[]}]`)
}

//
// Test creating a backtest
//
func TestCreateBacktests01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Create  a test screen
	screen := models.Screener{Name: "Test Screener", UserId: uint(1)}
	db.Create(&screen)

	// Create a test backtest
	btM := models.Backtest{
		UserId:          uint(2),
		StartingBalance: 5000.00,
		PositionSize:    "15-percent",
		StartDate:       models.Date{helpers.ParseDateNoError("2022-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2022-12-31")},
		Midpoint:        true,
		TradeSelect:     "highest-credit",
		Benchmark:       "SPY",
		Screen:          screen,
	}

	// Convert to json.
	j, _ := json.Marshal(btM)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/backtest", bytes.NewBuffer(j))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/backtest", c.CreateBacktest)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Backtest{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Get the backtest we just put in the db.
	b, err := db.BacktestGetById(1)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Id, b.Id)
	st.Expect(t, btM.StartingBalance, b.EndingBalance)
	st.Expect(t, uint(1), b.UserId)
	st.Expect(t, "SPY", b.Benchmark)
	st.Expect(t, uint(1), b.Screen.UserId)
	st.Expect(t, result.Id, b.Screen.BacktestId)
}

//
// Test creating a backtest - validation errors
//
func TestCreateBacktests02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Create a test backtest
	btM := models.Backtest{
		UserId:          uint(2),
		StartingBalance: 5000.00,
		EndingBalance:   5000.00,
		StartDate:       models.Date{helpers.ParseDateNoError("2022-01-01")},
		EndDate:         models.Date{helpers.ParseDateNoError("2022-12-31")},
		Midpoint:        true,
		Benchmark:       "SPY",
	}

	// Convert to json.
	j, _ := json.Marshal(btM)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/backtest", bytes.NewBuffer(j))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/backtest", c.CreateBacktest)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	st.Expect(t, w.Body.String(), `{"errors":{"position_size":"The position size field is required.","screen":"Screener Id is required.","trade_select":"The trade select field is required."}}`)
}
