//
// Date: 9/14/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"golang.org/x/crypto/bcrypt"
)

//
// TestGetProfile01
//
func TestGetProfile01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Set User
	db.Exec("TRUNCATE TABLE users;")
	db.Exec("TRUNCATE TABLE brokers;")
	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/me/profile", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/me/profile", c.GetProfile)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.User{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Email, "spicer+robtester@options.cafe")

	// Testing the json is best because we want to make sure something bad does not sneak in like a password.
	st.Expect(t, w.Body.String(), `{"id":1,"first_name":"Rob","last_name":"Tester","email":"spicer+robtester@options.cafe","phone":"","address":"","city":"","state":"","zip":"","country":"","brokers":[],"google_sub_id":"","last_activity":"0001-01-01T00:00:00Z"}`)
}

//
// TestUpdateProfile01
//
func TestUpdateProfile01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Set User
	db.Exec("TRUNCATE TABLE users;")
	db.Exec("TRUNCATE TABLE brokers;")
	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})

	// Body data
	var bodyStr = []byte(`{"first_name":"Mike","last_name":"Tester","email":"spicer+unittest@options.cafe","phone":"555-234-1234","address":"901 Brutscher Street, D112","city":"Newberg","state":"OR","zip":"97132","country":"USA"}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/me/profile", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/me/profile", c.UpdateProfile)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.User{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 202)
	st.Expect(t, result.Email, "spicer+unittest@options.cafe")

	// Testing the json is best because we want to make sure something bad does not sneak in like a password.
	st.Expect(t, w.Body.String(), `{"id":1,"first_name":"Mike","last_name":"Tester","email":"spicer+unittest@options.cafe","phone":"555-234-1234","address":"901 Brutscher Street, D112","city":"Newberg","state":"OR","zip":"97132","country":"USA","brokers":[],"google_sub_id":"","last_activity":"0001-01-01T00:00:00Z"}`)
}

//
// TestResetPassword01
//
func TestResetPassword01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Set User
	db.Exec("TRUNCATE TABLE users;")
	db.Exec("TRUNCATE TABLE brokers;")
	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active", Password: "$2a$10$eJ4biSke/V5Id9DK1nb2ZeCrGjI2IMaSQ.vTpaDeRbo4kg77RdhiC"})

	// Body data
	var bodyStr = []byte(`{"current_password":"foobar","new_password":"abc123"}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/me/rest-password", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/me/rest-password", c.ResetPassword)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 202)

	// Validate new password
	user, _ := db.GetUserById(uint(1))
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("abc123"))
	st.Expect(t, err, nil)

}

//
// TestResetPassword02
//
func TestResetPassword02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Set User
	db.Exec("TRUNCATE TABLE users;")
	db.Exec("TRUNCATE TABLE brokers;")
	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active", Password: "$2a$10$eJ4biSke/V5Id9DK1nb2ZeCrGjI2IMaSQ.vTpaDeRbo4kg77RdhiC"})

	// Body data
	var bodyStr = []byte(`{"current_password":"foobar","new_password":"abc"}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/me/rest-password", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/me/rest-password", c.ResetPassword)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Please enter a password at least 6 chars long."}`)
}

//
// TestResetPassword03
//
func TestResetPassword03(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Set User
	db.Exec("TRUNCATE TABLE users;")
	db.Exec("TRUNCATE TABLE brokers;")
	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active", Password: "$2a$10$eJ4biSke/V5Id9DK1nb2ZeCrGjI2IMaSQ.vTpaDeRbo4kg77RdhiC"})

	// Body data
	var bodyStr = []byte(`{"current_password":"foobar!!!","new_password":"abc123"}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/me/rest-password", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/me/rest-password", c.ResetPassword)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Incorrect current password."}`)
}

//
// TestUpdateCreditCard01
//
func TestUpdateCreditCard01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Set User
	db.Exec("TRUNCATE TABLE users;")

	// Create a test user.
	user := models.User{
		FirstName: "Jane",
		LastName:  "Unittester",
		Email:     "jane+unittest@options.cafe",
		Status:    "Trial",
	}

	// Insert user (create at stripe)
	err := db.CreateNewUserWithStripe(user)
	st.Expect(t, err, nil)

	// Since we are testing we know the user is id 1
	dbUser, err := db.GetUserById(1)
	st.Expect(t, err, nil)

	// Add a a credit card for testing we want to verify stripe only has one card at a time.
	err = db.UpdateCreditCard(dbUser, "tok_amex")
	st.Expect(t, err, nil)

	// Body data
	var bodyStr = []byte(`{"token":"tok_mastercard","coupon":"abc123"}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/me/update-credit-card", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/me/update-credit-card", c.UpdateCreditCard)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 202)

	// Get subscription with stripe
	sub, err := db.GetSubscriptionWithStripe(dbUser)
	st.Expect(t, err, nil)

	// All good?
	st.Expect(t, sub.CardBrand, "MasterCard")
	st.Expect(t, sub.CardLast4, "4444")
	st.Expect(t, sub.CardExpMonth, 10)
	st.Expect(t, sub.CardExpYear, 2019)

	// Clean things up at stripe
	db.DeleteUserWithStripe(dbUser)
}

/* End File */
