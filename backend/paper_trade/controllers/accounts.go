//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-18
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/cloudmanic/app.options.cafe/backend/paper_trade/models"
	"github.com/gin-gonic/gin"
)

//
// Return information me.
//
func (t *Controller) GetAccountMe(c *gin.Context) {

	// Make sure the AccountId.
	accountId := c.MustGet("accountId").(uint)

	// Get the account
	account := models.Account{}
	t.DB.New().Where("id = ?", accountId).Find(&account)

	// Return happy JSON
	c.JSON(200, account)
}

//
// Create a new paper trading account.
//
func (t *Controller) CreateAccount(c *gin.Context) {

	// t.DB.New().Save(&models.Account{
	//   FirstName:     "Spicer",
	//   LastName:      "Last",
	//   Email:         "spicer@options.cafe",
	//   AccountNumber: "abc123",
	//   AccessToken:   "asdfasdf34234234",
	// })

	// Return happy JSON
	c.JSON(200, gin.H{"status": "ok"})
}

/* End File */
