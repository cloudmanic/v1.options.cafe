//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

//
// Add a notification channel for this user.
//
func (t *Controller) CreateNotifyChannel(c *gin.Context) {

	// Setup NotifyChannel obj
	o := models.NotifyChannel{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o) != nil {
		return
	}

	// Make sure the UserId is correct.
	o.UserId = c.MustGet("userId").(uint)

	// See if we already have this entry by channelID
	s := models.NotifyChannel{}
	t.DB.New().Where("user_id = ? AND type = ? AND channel_id = ?", o.UserId, o.Type, o.ChannelId).Find(&s)

	// If we already have an obj return that one.
	if s.Id > 0 {
		t.RespondCreated(c, s, nil)
	} else {
		// Create Screen
		err := t.DB.CreateNewRecord(&o, models.InsertParam{})
		t.RespondCreated(c, o, err)
	}
}

/* End File */
