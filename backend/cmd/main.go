//
// Date: 11/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cmd

import (
	"flag"
	"fmt"

	"github.com/app.options.cafe/backend/library/import/options"
)

//
// Run this and see if we have any commands to run.
//
func Run() bool {

	// Grab flags
	action := flag.String("cmd", "none", "bulk-eod-options-import")
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

	// Just a test
	case "test":
		fmt.Println("CMD Works....")
		return true
		break

	}

	// This is not a command. Run the server.
	return false
}

/* End File */
