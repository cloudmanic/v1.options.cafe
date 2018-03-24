//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cron

import (
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/cron/data_import"
	"github.com/cloudmanic/app.options.cafe/backend/library/import/options"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/jasonlvhit/gocron"
)

//
// Start....
//
func Start(db *models.DB) {

	// Setup instance
	d := data_import.Base{DB: db}

	// Lets get started
	services.Critical("Cron Started: " + os.Getenv("APP_ENV"))

	// Setup jobs we need to run
	gocron.Every(1).Day().At("14:00").Do(d.DoSymbolImport)
	gocron.Every(1).Day().At("22:00").Do(options.DoEodOptionsImport)

	// function Start start all the pending jobs
	<-gocron.Start()

}

/* End File */
