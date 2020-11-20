//
// Date: 2018-11-11
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-11
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package worker

import "app.options.cafe/models"

type JobRequest struct {
	DB       models.Datastore
	Action   string `json:"action"`
	UserId   uint   `json:"user_id"`
	BrokerId uint   `json:"broker_id"`
}

/* End File */
