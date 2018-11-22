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

	// Stuff we do on start as well
	user.ExpireTrails(db)
	user.ClearExpiredSessions(db)

	// Lets get started
	services.Critical("Cron Started: " + os.Getenv("APP_ENV"))

	// New Cron instance
	c := cron.New()

	// Setup jobs we need to run
	c.AddFunc("0 0 14 * * *", func() { data_import.DoSymbolImport(db) }) // Every day at 14:00
	c.AddFunc("0 0 22 * * *", func() { data_import.DoSymbolImport(db) }) // Every day at 22:00

	// User clean up stuff
	c.AddFunc("@hourly", func() { user.ExpireTrails(db) })
	c.AddFunc("@every 12h", func() { user.ClearExpiredSessions(db) })

	// System stuff.
	c.AddFunc("@every 10s", func() { DatabasePing(db) })

	// Start cron service
	c.Run()
}

//
// We use this to keep the database alive.
//
func DatabasePing(db *models.DB) {

	// Just run a query to make sure things are active.
	a := []models.Application{}
	db.New().Find(&a)

}

/* End File */
