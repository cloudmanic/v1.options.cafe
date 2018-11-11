//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cron

import (
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/cron/data_import"
	"github.com/cloudmanic/app.options.cafe/backend/cron/user"
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
	u := user.Base{DB: db}
	d := data_import.Base{DB: db}

	// Stuff we do on start as well
	u.ExpireTrails()
	u.ClearExpiredSessions()

	// Lets get started
	services.Critical("Cron Started: " + os.Getenv("APP_ENV"))

	// Setup jobs we need to run
	gocron.Every(1).Day().At("14:00").Do(d.DoSymbolImport)
	gocron.Every(1).Day().At("22:00").Do(options.DoEodOptionsImport)

	// User clean up stuff
	gocron.Every(1).Hour().Do(u.ExpireTrails)
	gocron.Every(12).Hours().Do(u.ClearExpiredSessions)

	// function Start start all the pending jobs
	<-gocron.Start()

}

/* End File */
