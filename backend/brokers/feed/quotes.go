//
// Date: 2/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import (
	"fmt"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

//
// Ticker - Get DetailedQuotes : 1 second
//
func (t *Base) DoGetDetailedQuotes() {

	for {

		// Load up our DetailedQuotes
		err := t.GetActiveSymbolsDetailedQuotes()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 1 second
		time.Sleep(time.Second * 1)

	}

}

//
// Do get quotes - more details from the streaming - activeSymbols
//
func (t *Base) GetActiveSymbolsDetailedQuotes() error {

	symbols := t.GetActiveSymbols()
	detailedQuotes, err := t.Api.GetQuotes(symbols)

	if err != nil {
		return err
	}

	// Loop through the quotes sending them up the websocket channel
	for _, row := range detailedQuotes {

		// Send up websocket.
		err = t.WriteDataChannel("quote", row)

		if err != nil {
			return fmt.Errorf("GetActiveSymbolsDetailedQuotes() WriteDataChannel : ", err)
		}

	}

	// Return happy
	return nil

}

/* End File */
