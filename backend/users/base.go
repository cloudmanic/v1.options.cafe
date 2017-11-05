//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import (
	"app.options.cafe/backend/controllers"
	"app.options.cafe/backend/models"
)

type Base struct {
	DB              *models.DB
	Users           map[uint]*User
	DataChan        chan controllers.SendStruct
	QuoteChan       chan controllers.SendStruct
	FeedRequestChan chan controllers.SendStruct
}

/* End File */
