//
// Date: 9/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

//
// Ticker - Positions : 1 hour
//
func (t *Base) DoPositionsTicker() {

	var err error

	for {

		// Do we break out ?
		t.MuPolling.Lock()
		breakOut := t.Polling
		t.MuPolling.Unlock()

		if !breakOut {
			break
		}

		err = t.GetPositions()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 1 hour
		time.Sleep(time.Hour * 1)

	}

	services.Info("Stopping DoPositionsTicker() : " + t.User.Email)
}

//
// Fetch - Positions
//
func (t *Base) GetPositions() error {

	// Make API call
	positions, err := t.Api.GetPositions()

	if err != nil {
		return err
	}

	// Loop through and add any positions to the active_symbols table.
	for _, row := range positions {
		t.DB.CreateActiveSymbol(t.User.Id, row.Symbol)
	}

	// Return Happy
	return nil
}

/* End File */
