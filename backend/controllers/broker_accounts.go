//
// Date: 9/21/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

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

	// Return success.
	c.JSON(http.StatusNoContent, nil)
}

/* End File */
