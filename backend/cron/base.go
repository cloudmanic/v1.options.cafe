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
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/robfig/cron"
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

	// New Cron instance
	c := cron.New()

	// Setup jobs we need to run
	c.AddFunc("0 0 14 * * *", d.DoSymbolImport) // Every day at 14:00
	c.AddFunc("0 0 22 * * *", d.DoSymbolImport) // Every day at 22:00

	// User clean up stuff
	c.AddFunc("@hourly", u.ExpireTrails)
	c.AddFunc("@every 12h", u.ClearExpiredSessions)

	// Start cron service
	c.Run()
}

/* End File */
