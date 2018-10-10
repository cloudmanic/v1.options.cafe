//
// Date: 9/14/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
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

/* End File */
