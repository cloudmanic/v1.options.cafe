//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"go/build"

	"github.com/cloudmanic/app.options.cafe/backend/paper_trade/models"
	env "github.com/jpfuentes2/go-env"
)

type Controller struct {
	DB models.Datastore
}

//
// Start up the controller.
//
func init() {
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/paper_trade/.env")
}

/* End File */
