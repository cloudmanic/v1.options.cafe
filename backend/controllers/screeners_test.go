//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestGetScreeners01
//
func TestGetScreeners01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

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
	db, _ := models.NewDB()

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
	db, _ := models.NewDB()

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
	db, _ := models.NewDB()

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
