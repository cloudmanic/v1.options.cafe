//
// Date: 2/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

//
// Ticker - See if we need to refresh an access token : 60 seconds
//
func (t *Base) DoAccessTokenRefresh() {
	var err error

	for {

		err = t.AccessTokenRefresh()

		if err != nil {
			services.Warning(err)
		}

		// Sleep for 60 second.
		time.Sleep(time.Second * 60)

	}
}

//
// Do update access token from refresh
//
func (t *Base) AccessTokenRefresh() error {
	err := t.Api.DoRefreshAccessTokenIfNeeded(t.User)

	if err != nil {
		return err
	}

	// Return Happy
	return nil
}

/* End File */
