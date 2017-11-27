//
// Date: 11/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"flag"
	"runtime"

	"github.com/app.options.cafe/backend/library/import/options"
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

	// Grab flags
	action := flag.String("action", "none", "bulk-eod-options-import")
	flag.Parse()

	switch *action {

	// Symbol import
	case "bulk-eod-options-import":
		options.DoBulkEodImportToPerSymbolDay()
		break

	}

}

/* End File */
