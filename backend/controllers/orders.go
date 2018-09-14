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
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

//
// Preview an order
//
func (t *Controller) PreviewOrder(c *gin.Context) {

	// Build request
	order, brokerCont, brokerAccount, err := orderBuildRequest(t, c)

	if err != nil {
		services.Warning(err)
		return
	}

	// Send preview request to broker
	preview, err := brokerCont.PreviewOrder(brokerAccount.AccountNumber, order)

	if err != nil {
		services.Warning(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Preview error."})
		return
	}

	// Return happy JSON
	c.JSON(200, preview)
}

//
// Place an order
//
func (t *Controller) SubmitOrder(c *gin.Context) {

	// Build request
	order, brokerCont, brokerAccount, err := orderBuildRequest(t, c)

	if err != nil {
		services.Warning(err)
		return
	}

	// Send order request to broker
	orderRequest, err := brokerCont.SubmitOrder(brokerAccount.AccountNumber, order)

	if err != nil {
		services.Warning(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order error."})
		return
	}

	// Return happy JSON
	c.JSON(200, orderRequest)
}

// --------------- Helper Function --------------------- //

//
// Shared logic between
//
func orderBuildRequest(t *Controller, c *gin.Context) (types.Order, tradier.Api, models.BrokerAccount, error) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	var order types.Order

	// Json to Object
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return types.Order{}, tradier.Api{}, models.BrokerAccount{}, err
	}

	// Get broker account.
	brokerAccount, err := t.DB.GetBrokerAccountByIdUserId(order.BrokerAccountId, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown broker. (#001)"})
		return types.Order{}, tradier.Api{}, models.BrokerAccount{}, err
	}

	// Get the broker
	broker, err := t.DB.GetBrokerById(brokerAccount.BrokerId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown broker. (#002)"})
		return types.Order{}, tradier.Api{}, models.BrokerAccount{}, err
	}

	// Decrypt the access token
	apiKey, err := helpers.Decrypt(broker.AccessToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown broker. (#003)"})
		return types.Order{}, tradier.Api{}, models.BrokerAccount{}, err
	}

	var brokerCont tradier.Api

	// Figure out which broker connection to setup.
	switch broker.Name {

	case "Tradier":
		brokerCont = tradier.Api{ApiKey: apiKey, DB: t.DB, Sandbox: false}

	case "Tradier Sandbox":
		brokerCont = tradier.Api{ApiKey: apiKey, DB: t.DB, Sandbox: true}

	default:
		services.Critical("Order: Unknown Broker : " + broker.Name)

	}

	// Return happy
	return order, brokerCont, brokerAccount, nil
}

/* End File */
