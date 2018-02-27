//
// Date: 2/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/library/state"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/websocket"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Test - GetMarketStatus
//
func TestGetMarketStatus01(t *testing.T) {

	// Create test status
	status := websocket.MarketStatus{
		Date:        "2017-09-04",
		State:       "closed",
		Description: "Market is closed.",
	}

	// Store in our caching system
	state.SetMarketStatus(status)

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/status/market", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/status/market", c.GetMarketStatus)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := websocket.MarketStatus{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, result.Date, "2017-09-04")
	st.Expect(t, result.State, "closed")
	st.Expect(t, result.Description, "Market is closed.")
}

/* End File */
