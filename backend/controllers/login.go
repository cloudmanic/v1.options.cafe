//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/app.options.cafe/backend/library/realip"
	"github.com/app.options.cafe/backend/library/services"
	"github.com/gin-gonic/gin"
)

// TODO: Lots of duplicate code in here with setting headers and such. Should clean up. Also see Forgot Password, and Register.

//
// Login to account.
//
func (t *Controller) DoLogin(c *gin.Context) {

	// Set response
	if os.Getenv("APP_ENV") == "local" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// Decode json passed in
	decoder := json.NewDecoder(c.Request.Body)

	type LoginPost struct {
		Email    string
		Password string
	}

	var post LoginPost

	err := decoder.Decode(&post)

	if err != nil {
		services.Error(err, "DoLogin - Failed to decode JSON posted in")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."})
		return
	}

	defer c.Request.Body.Close()

	// Validate user.
	if err := t.DB.ValidateUserLogin(post.Email, post.Password); err != nil {

		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Login user in by email and password
	user, err := t.DB.LoginUserByEmailPass(post.Email, post.Password, c.Request.UserAgent(), realip.RealIP(c.Request))

	if err != nil {
		services.Error(err, "DoLogin - Unable to log user in. (CreateUser)")

		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, we could not find your account."})
		return
	}

	// Here we check to see if we have any brokers. If there are no brokers the user needs to select at least one to do anything.
	var brokerCount = len(user.Brokers)

	type Response struct {
		Status      uint   `json:"status"`
		UserId      uint   `json:"user_id"`
		AccessToken string `json:"access_token"`
		BrokerCount int    `json:"broker_count"`
	}

	resObj := &Response{
		Status:      1,
		UserId:      user.Id,
		AccessToken: user.Session.AccessToken,
		BrokerCount: brokerCount,
	}

	// Return success json.
	c.JSON(200, resObj)
}

/* End File */
