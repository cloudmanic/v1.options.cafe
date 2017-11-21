//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/app.options.cafe/backend/cron/data_import"
	"github.com/app.options.cafe/backend/library/services"
	"github.com/app.options.cafe/backend/models"
	"github.com/jasonlvhit/gocron"
	env "github.com/jpfuentes2/go-env"
)

//
// Main....
//
func main() {

	// Setup CPU stuff.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load ENV (if we have it.)
	env.ReadEnv("../.env")

	// Start the db connection.
	db, err := models.NewDB()

	if err != nil {
		services.Fatal(err)
	}

	// Close db when this app dies. (This might be useless)
	defer db.Close()

	// Grab flags
	action := flag.String("action", "none", "If you want to run just one command you can use action. { symbol-import }")
	flag.Parse()

	// Setup instance
	d := data_import.Base{DB: db}

	// Run one action at a time or start cron.
	if *action != "none" {

		switch *action {

		// Symbol import
		case "symbol-import":
			d.DoSymbolImport()
			break

		}

	} else {

		// Lets get started
		services.Critical("Cron Started: " + os.Getenv("APP_ENV"))

		// Setup jobs we need to run
		gocron.Every(1).Day().At("14:00").Do(d.DoSymbolImport)
		gocron.Every(1).Day().At("22:00").Do(data_import.DoEodOptionsImport)

		// function Start start all the pending jobs
		<-gocron.Start()
	}

}

/* End File */
