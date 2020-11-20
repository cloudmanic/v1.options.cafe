//
// Date: 2018-11-16
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-17
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app.options.cafe/library/analyze"
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Test - AnalyzeOptionsProfitLossByUnderlyingPrice - 01 - Success
//
func TestAnalyzeOptionsProfitLossByUnderlyingPrice01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Setup Post STring
	postStr := `{
	"open_cost": 157.00,
	"legs": [
		{ "symbol_str": "SPY181221C00250000", "qty": 1 },
		{ "symbol_str": "SPY181221C00260000", "qty": -2 },
		{ "symbol_str": "SPY181221C00270000", "qty": 1 }
	]}`

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/analyze/options/underlying-price", bytes.NewBuffer([]byte(postStr)))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/analyze/options/underlying-price", c.AnalyzeOptionsProfitLossByUnderlyingPrice)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var results []analyze.Result

	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Validate result
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, len(results), 501)
	st.Expect(t, results[100].Profit, -157.00)
	st.Expect(t, helpers.Round(results[100].UnderlyingPrice, 2), 246.70)
	st.Expect(t, results[300].Profit, 333.00)
	st.Expect(t, helpers.Round(results[300].UnderlyingPrice, 2), 265.10)
}

/* End File */
