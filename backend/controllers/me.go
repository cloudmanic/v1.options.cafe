//
// Date: 9/14/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/gin-gonic/gin"
)

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
