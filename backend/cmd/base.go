//
// Date: 11/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cmd

import (
	"flag"
	"fmt"

	"github.com/cloudmanic/app.options.cafe/backend/cmd/actions"
	"github.com/cloudmanic/app.options.cafe/backend/cron"
	"github.com/cloudmanic/app.options.cafe/backend/library/backtest"
	"github.com/cloudmanic/app.options.cafe/backend/library/import/options"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Run this and see if we have any commands to run.
//
func Run(db *models.DB) bool {

	// Grab flags
	action := flag.String("cmd", "none", "")
	name := flag.String("name", "", "")
	symbol := flag.String("symbol", "", "")
	flag.Parse()

	switch *action {

	// Bulk EOD options import
	case "bulk-eod-options-import":
		options.DoBulkEodImportToPerSymbolDay()
		return true
		break

	// EOD options import
	case "eod-options-import":
		options.DoEodOptionsImport()
		return true
		break

	// Download a EOD options symbol locally (go run main.go -cmd=download-eod-symbol -symbol=spy)
	case "download-eod-symbol":
		backtest.DownloadEodSymbol(*symbol, true)
		return true
		break

	// Create a new application from the CLI
	case "create-application":
		actions.CreateApplication(db, *name)
		return true
		break

	// Start cron server
	case "cron":
		cron.Start(db)
		return true
		break

	// Just a test
	case "test":
		fmt.Println("CMD Works....")
		return true
		break

	}

	return false
}

/* End File */
