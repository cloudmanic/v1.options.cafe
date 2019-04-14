//
// Date: 9/21/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"errors"
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/queue"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

//
// Get the balance of a particular broker account.
//
func (t *Controller) BrokerAccountGetBalance(c *gin.Context) {
	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Set as int
	brokerId, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Set as int
	acctId, err := strconv.ParseInt(c.Param("acctId"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the broker account by id.
	brokerAccount, err := t.DB.GetBrokerAccountByIdUserId(uint(acctId), userId)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Make sure we have a valid broker Id
	if uint(brokerId) != brokerAccount.BrokerId {
		t.RespondError(c, errors.New("Broker not found."), httpGenericErrMsg)
		return
	}

	// Get the broker
	broker, err := t.DB.GetBrokerById(uint(brokerId))

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Decrypt the access token
	apiKey, err := helpers.Decrypt(broker.AccessToken)

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	var brokerCont tradier.Api

	// Figure out which broker connection to setup.
	switch broker.Name {

	case "Tradier":
		brokerCont = tradier.Api{ApiKey: apiKey, DB: t.DB, Sandbox: false}

	case "Tradier Sandbox":
		brokerCont = tradier.Api{ApiKey: apiKey, DB: t.DB, Sandbox: true}

	default:
		services.InfoMsg("Order: Unknown Broker : " + broker.Name)

	}

	// Make API call to broker.
	result, err := brokerCont.GetBalances()

	if err != nil {
		t.RespondError(c, err, httpGenericErrMsg)
		return
	}

	// Just find the result we want to return
	for _, row := range result {

		if row.AccountNumber == brokerAccount.AccountNumber {

			// Return happy
			c.JSON(200, row)
			return

		}

	}

	// Return nothing found JSON
	t.RespondError(c, errors.New("Broker account not found."), httpGenericErrMsg)
}

//
// Update  broker account. We only allow updates as the oauth process
// manages the accounts mostly.
//
func (t *Controller) UpdateBrokerAccount(c *gin.Context) {

	// Setup BrokerAccount obj
	o := models.BrokerAccount{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// Make sure the UserId is correct.
	o.UserId = c.MustGet("userId").(uint)

	// Set as int
	brokerId, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Set as int
	acctId, err := strconv.ParseInt(c.Param("acctId"), 10, 32)

	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}

	// Get the broker account by id.
	orgObj, err := t.DB.GetBrokerAccountByIdUserId(uint(acctId), o.UserId)

	if t.RespondError(c, err, httpNoRecordFound) {
		return
	}

	// Make sure this broker is connected
	if orgObj.BrokerId != uint(brokerId) {
		t.RespondError(c, errors.New("Unknown Broker account."), httpNoRecordFound)
		return
	}

	// Set id.
	o.Id = orgObj.Id

	// Maintain broker id.
	o.BrokerId = orgObj.BrokerId

	// We override account number. That is not something a user should change.
	o.AccountNumber = orgObj.AccountNumber

	// Set the create date.
	o.CreatedAt = orgObj.CreatedAt

	// Update BrokerAccount
	t.DB.New().Save(&o)

	// Send websocket with broker change
	queue.Write("oc-websocket-write", `{"uri":"change-detected","user_id":`+strconv.Itoa(int(o.UserId))+`,"body": { "type": "brokers" } }`)

	// Return success.
	c.JSON(204, o)
}

/* End File */
