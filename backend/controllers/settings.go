//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

//
// Get a settings by UserId.
//
func (t *Controller) GetSettings(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Get the settings by id.
	settings := t.DB.SettingsGetOrCreateByUserId(userId)

	// Return happy JSON
	c.JSON(200, settings)
}

//
// Update a settings by UserId.
//
func (t *Controller) UpdateSettings(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Get the settings by id.
	settings := t.DB.SettingsGetOrCreateByUserId(userId)

	// Setup Settings obj
	o := models.Settings{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// We only allow a few fields to be updated via the API
	settings.StrategyPcsClosePrice = o.StrategyPcsClosePrice
	settings.StrategyPcsOpenPrice = o.StrategyPcsOpenPrice
	settings.StrategyPcsLots = o.StrategyPcsLots
	settings.StrategyCcsClosePrice = o.StrategyCcsClosePrice
	settings.StrategyCcsOpenPrice = o.StrategyCcsOpenPrice
	settings.StrategyCcsLots = o.StrategyCcsLots
	settings.StrategyPdsClosePrice = o.StrategyPdsClosePrice
	settings.StrategyPdsOpenPrice = o.StrategyPdsOpenPrice
	settings.StrategyPdsLots = o.StrategyPdsLots
	settings.StrategyCdsClosePrice = o.StrategyCdsClosePrice
	settings.StrategyCdsOpenPrice = o.StrategyCdsOpenPrice
	settings.StrategyCdsLots = o.StrategyCdsLots
	settings.NoticeTradeFilledEmail = o.NoticeTradeFilledEmail
	settings.NoticeTradeFilledSms = o.NoticeTradeFilledSms
	settings.NoticeTradeFilledPush = o.NoticeTradeFilledPush
	settings.NoticeMarketOpenedEmail = o.NoticeMarketOpenedEmail
	settings.NoticeMarketOpenedSms = o.NoticeMarketOpenedSms
	settings.NoticeMarketOpenedPush = o.NoticeMarketOpenedPush
	settings.NoticeMarketClosedEmail = o.NoticeMarketClosedEmail
	settings.NoticeMarketClosedSms = o.NoticeMarketClosedSms
	settings.NoticeMarketClosedPush = o.NoticeMarketClosedPush

	// Update settings
	t.DB.New().Save(&settings)

	// Return happy JSON
	c.JSON(http.StatusNoContent, gin.H{})
}

/* End File */
