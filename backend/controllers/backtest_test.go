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
