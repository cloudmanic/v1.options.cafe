//
// Date: 6/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Test - DoLogOut 01
//
func TestDoLogOut01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Add test session
	db.Create(&models.Session{AccessToken: "abc123456789you"})

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/oauth/logout?access_token=abc123456789you", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.GET("/oauth/logout", c.DoLogOut)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Test result
	st.Expect(t, w.Body.String(), `{"status":"ok"}`)
}

//
// Test - DoLogOut 02
//
func TestDoLogOut02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/oauth/logout?access_token=notfound", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.GET("/oauth/logout", c.DoLogOut)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Test result
	st.Expect(t, w.Body.String(), `{"error":"Sorry, we could not find your session.","status":"error"}`)
}

//
// Test - DoLogOut 03
//
func TestDoLogOut03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/oauth/logout", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.GET("/oauth/logout", c.DoLogOut)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Test result
	st.Expect(t, w.Body.String(), `{"error":"Sorry, access_token is required."}`)
}

/* End File */
