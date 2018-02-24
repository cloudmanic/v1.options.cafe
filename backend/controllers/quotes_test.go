//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

//
// Test - GetHistoricalQuotes
//
func TestGetHistoricalQuotes01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/quotes/historical", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()
	r.GET("/api/v1/quotes/historical", c.GetHistoricalQuotes)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	spew.Dump(w.Body.String())

	// blogs := []models.Blog{}

	// // This is how you ALWAYS convert from a string to byte array.
	// err := json.Unmarshal([]byte(w.Body.String()), &blogs)
	//spew.Dump(blogs)

	// Parse json that returned.
	//st.Expect(t, err, nil)
	// st.Expect(t, len(blogs), 1)
	// st.Expect(t, blogs[0].Body, "This is a blog Body....")
	// st.Expect(t, blogs[0].Title, "My First Blog Post")
	// This is how you would test date time.
	//st.Expect(t, blogs[0].PostDate, time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC))

}

/* End File */
