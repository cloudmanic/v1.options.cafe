//
// Date: 4/13/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

// Package seed used to load data data into our database.
package seed

import (
	"os"
	"os/exec"

	"app.options.cafe/library/services"
	"app.options.cafe/models"
)

//
// SeedLocalDatabase - This is used for local development. Just make sure the
// required data is in the database.
//
func LocalDatabase(db models.Datastore) {
	// Make sure we are local env
	if os.Getenv("APP_ENV") != "local" {
		services.InfoMsg("Skipping local seeding. Not in local dev mode.")
		return
	}

	// Make sure our db is empty
	count := 0
	db.New().Model(&models.Application{}).Count(&count)

	// If we have at least one record we assume we already seeded the DB.
	if count > 0 {
		services.InfoMsg("Skipping local seeding. Already seeded.")
		return
	}

	// Test as a db import functionality so might as well use it.
	LoadSqlDump(db, "symbols")
	LoadSqlDump(db, "applications")
	LoadSqlDump(db, "historical_quotes")
}

//
// LoadSqlDump - Take a dataset string and load content from
// a data file into our database. These dataset strings are mysql dumps
//
func LoadSqlDump(db models.Datastore, dataSet string) {
	// Set file
	dataFile := os.Getenv("CODE_LOCATION") + "/library/seed/data/" + dataSet + ".sql"

	// Set CMD
	cmd := "mysql --host=127.0.0.1 --port=" + os.Getenv("DB_PORT") + " -u " + os.Getenv("DB_USERNAME") + " -p" + os.Getenv("DB_PASSWORD") + " " + os.Getenv("DB_DATABASE") + " < " + dataFile

	// Run CLI command to import data
	_, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		services.Info(err)
	}

	services.InfoMsg("Seeding the " + dataSet + " table.")
}
