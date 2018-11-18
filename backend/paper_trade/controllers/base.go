//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-18
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"go/build"

	masterModels "github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/paper_trade/models"
	env "github.com/jpfuentes2/go-env"
)

var masterDB masterModels.Datastore

type Controller struct {
	DB models.Datastore
}

//
// Start up the controller.
//
func init() {
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/paper_trade/.env")

	// Start the db connection to our master db.
	masterDB, _ := masterModels.NewDB()

	// Ping DB just so we do not have compile errors
	masterDB.New().DB().Ping()
}

/* End File */
