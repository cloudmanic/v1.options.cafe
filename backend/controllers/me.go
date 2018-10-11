//
// Date: 9/14/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"
)

//
// Get Me. Return user profile.
//
func (t *Controller) GetProfile(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Return happy JSON
	c.JSON(200, user)
}

//
// Update Me. Update a user profile.
//
func (t *Controller) UpdateProfile(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Setup BrokerAccount obj
	o := models.User{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// We only allow a few fields to be updated via the API
	user.Email = o.Email
	user.FirstName = o.FirstName
	user.LastName = o.LastName
	user.Phone = o.Phone
	user.Address = o.Address
	user.City = o.City
	user.State = o.State
	user.Country = o.Country
	user.Zip = o.Zip

	// Update BrokerAccount
	t.DB.New().Save(&user)

	// Return happy JSON
	c.JSON(202, user)
}

//
// Get the current subscription details.
//
func (t *Controller) GetSubscription(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Get subscription with stripe
	sub, err := t.DB.GetSubscriptionWithStripe(user)

	if t.RespondError(c, err, "Subscription not found. Please contact help@options.cafe") {
		return
	}

	// Return happy JSON
	c.JSON(200, sub)
}

//
// Reset the users password
//
func (t *Controller) ResetPassword(c *gin.Context) {

	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.DB.GetUserById(userId)

	if t.RespondError(c, err, "User not found. Please contact help@options.cafe") {
		return
	}

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	newPass := gjson.Get(string(body), "new_password").String()
	currentPass := gjson.Get(string(body), "current_password").String()

	// Now that we know the user lets make sure the password that was posted in was at least 6 chars.
	err = t.DB.ValidatePassword(newPass)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a password at least 6 chars long."})
		return
	}

	// Validate password here by comparing hashes nil means success
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPass))

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect current password."})
		return
	}

	// Change password
	err = t.DB.ResetUserPassword(user.Id, newPass)

	if err != nil {
		services.BetterError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to reset password."})
		return
	}

	// Return happy
	c.JSON(202, nil)
}

/* End File */
