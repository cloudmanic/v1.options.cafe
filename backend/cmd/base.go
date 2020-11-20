//
// Date: 11/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cmd

import (
	"flag"

	"app.options.cafe/cmd/actions"
	"app.options.cafe/cron"
	"app.options.cafe/cron/data_import"
	"app.options.cafe/library/import/options"
	"app.options.cafe/library/polling"
	"app.options.cafe/library/worker/jobs"
	"app.options.cafe/models"
)

//
// Run this and see if we have any commands to run.
//
func Run(db *models.DB) bool {

	// Grab flags
	userId := flag.Int("user_id", 0, "")
	brokerAccountId := flag.Int("broker_account_id", 0, "")
	action := flag.String("cmd", "none", "--cmd={action}")
	name := flag.String("name", "", "")
	flag.Parse()

	switch *action {

	// This is a worker
	case "worker":
		jobs.Start(db)
		return true
		break

	// This is a broker feed poller
	case "broker-feed-poller":
		polling.Start(db)
		return true
		break

	// Go to Angular
	case "go-to-angular":
		actions.GoToAngular()
		return true
		break

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

	// Create a new application from the CLI
	case "create-application":
		actions.CreateApplication(db, *name)
		return true
		break

	// Download all the symbols Tradier knows about. (typically run daily from cron)
	case "symbol-import":
		data_import.DoSymbolImport(db)
		return true
		break

	// Start cron server
	case "cron":
		cron.Start(db)
		return true
		break

	// Return the json a broker give us by a userid
	case "broker-json-history":
		actions.GetJsonBrokerHistory(db, *userId, *brokerAccountId)
		return true
		break

	// Rebuild trade groups by user.
	case "user-rebuild-trade-groups":
		actions.UserRebuildTradeGroups(db, *userId)
		return true
		break

	// Run a backtest
	case "backtest-run":
		actions.RunBackTest(db, *userId)
		return true
		break

	}

	return false
}

/* End File */
