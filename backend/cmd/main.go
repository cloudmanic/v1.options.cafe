//
// Date: 11/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"flag"
	"fmt"

	"github.com/app.options.cafe/backend/library/import/options"
)

//
// Run this and see if we have any commands to run.
//
func main() {

	// Grab flags
	action := flag.String("cmd", "none", "bulk-eod-options-import")
	flag.Parse()

	switch *action {

	// Bulk EOD options import
	case "bulk-eod-options-import":
		options.DoBulkEodImportToPerSymbolDay()
		return
		break

	// EOD options import
	case "eod-options-import":
		options.DoEodOptionsImport()
		return
		break

	// Just a test
	case "test":
		fmt.Println("CMD Works....")
		return
		break

	}
}

/* End File */
