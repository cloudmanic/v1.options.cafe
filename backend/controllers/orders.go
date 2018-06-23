//
// Date: 6/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/gin-gonic/gin"
)

//
// Preview an order
//
func (t *Controller) PreviewOrder(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	var order types.Order

	// Json to Object
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get broker account.
	brokerAccount, err := t.DB.GetBrokerAccountByIdUserId(order.BrokerAccountId, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown broker. (#001)"})
		return
	}

	// Get the broker
	broker, err := t.DB.GetBrokerById(brokerAccount.BrokerId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown broker. (#002)"})
		return
	}

	// Decrypt the access token
	apiKey, err := helpers.Decrypt(broker.AccessToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown broker. (#003)"})
		return
	}

	// Setup the broker
	brokerCont := tradier.Api{
		DB:     t.DB,
		ApiKey: apiKey,
	}

	// Send preview request to broker
	preview, err := brokerCont.PreviewOrder(brokerAccount.AccountNumber, order)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Preview error."})
		return
	}

	// Return happy JSON
	c.JSON(200, preview)
}

/* End File */
