//
// Date: 3/5/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

//
// Test - WatchlistAdd
//
func TestWatchlistAdd01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{ "symbol_id": 5 }`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/watchlists/1/add", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/watchlists/1/add", c.WatchlistAdd)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Println(string(w.Body.String()))

	// // Grab result and convert to strut
	// result := []models.Broker{}
	// err := json.Unmarshal([]byte(w.Body.String()), &result)

	// // Parse json that returned.
	// st.Expect(t, err, nil)
	// st.Expect(t, len(result), 3)
	// st.Expect(t, result[0].Id, uint(1))
	// st.Expect(t, result[1].Id, uint(2))
	// st.Expect(t, result[2].Id, uint(3))
	// st.Expect(t, result[0].Name, "Tradier")
	// st.Expect(t, result[1].Name, "Tradeking")
	// st.Expect(t, result[2].Name, "Etrade")
	// st.Expect(t, len(result[0].BrokerAccounts), 2)
	// st.Expect(t, len(result[1].BrokerAccounts), 0)
	// st.Expect(t, len(result[2].BrokerAccounts), 0)
	// st.Expect(t, result[0].BrokerAccounts[0].AccountNumber, "YYY123ZY")
	// st.Expect(t, result[0].BrokerAccounts[1].AccountNumber, "ABC123ZY")
	// st.Expect(t, result[0].BrokerAccounts[0].BrokerId, uint(1))
	// st.Expect(t, result[0].BrokerAccounts[1].BrokerId, uint(1))
	// st.Expect(t, result[0].BrokerAccounts[0].AccountNumber, "YYY123ZY")
	// st.Expect(t, result[0].BrokerAccounts[1].AccountNumber, "ABC123ZY")

}

/* End File */
