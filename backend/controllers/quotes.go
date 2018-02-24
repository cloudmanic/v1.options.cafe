//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/gin-gonic/gin"
)

//
// Get time sales data
//
// TODO: pull the access token from the user's account.
//
func (t *Controller) GetHistoricalQuotes(c *gin.Context) {

	// // Get the user. This should never error because of the middleware
	// user, err := t.DB.GetUserById(c.MustGet("userId").(uint))

	// if err != nil {
	// 	t.RespondError(c, err, httpGenericErrMsg)
	// 	return
	// }

	// Setup the broker
	broker := tradier.Api{
		DB:     t.DB,
		ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"),
	}

	// Make API call to broker.
	result, err := broker.GetBalances()

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Return happy JSON
	c.JSON(200, result)
}

/* End File */
