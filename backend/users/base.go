//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import (
	"github.com/cloudmanic/app.options.cafe/backend/brokers/feed"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/websocket"
)

type Base struct {
	DB          *models.DB
	Users       map[uint]*UserFeed
	WsWriteChan chan websocket.SendStruct
	ActionChan  chan UserFeedAction
}

type UserFeed struct {
	Profile     models.User
	WsWriteChan chan websocket.SendStruct
	BrokerFeed  map[uint]*feed.Base
}

type UserFeedAction struct {
	Action string
	UserId uint
}

/* End File */
