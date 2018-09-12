//
// Date: 9/11/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import (
	"fmt"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/archive"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

//
// Ticker - Get GetHistory: 24 hours
//
func (t *Base) DoGetHistoryTicker() {
	var err error

	for {

		// Get history from tradier
		err = t.GetHistory()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 24 hours
		time.Sleep(time.Hour * 24)
	}
}

//
// Do get History
//
func (t *Base) GetHistory() error {
	history, err := t.Api.GetAllHistory()

	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("GetHistory(): ", err)
	}

	// Store the history in our database
	err = archive.StoreBrokerEvents(t.DB, history, t.User.Id, t.BrokerId)

	// Return Happy
	return nil
}

/* End File */
