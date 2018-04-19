//
// Date: 4/19/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Test AddActiveSymbol - 01
//
func TestAddActiveSymbol01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{"symbol":"SPY180525P00257000"}`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/symbols/add-active-symbol", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/symbols/add-active-symbol", c.AddActiveSymbol)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Body.String(), `{"id":2,"symbol":"SPY180525P00257000"}`)
}

//
// Test AddActiveSymbol - 02 (error)
//
func TestAddActiveSymbol02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	var postStr = []byte(`{"symbol":""}`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/symbols/add-active-symbol", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/symbols/add-active-symbol", c.AddActiveSymbol)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Body.String(), `{"error":"symbol field can not be empty."}`)
}

/* End File */
