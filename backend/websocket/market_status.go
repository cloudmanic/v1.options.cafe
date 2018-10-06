//
// Date: 2/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/notify/sms_push"
	"github.com/cloudmanic/app.options.cafe/backend/library/notify/web_push"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cnf/structhash"
)

type MarketStatus struct {
	Date        string `json:"date"`
	State       string `json:"state"`
	Description string `json:"description"`
}

//
// Monitor the market status via the Tradier
// admin account and send updates to all connected clients.
//
func (t *Controller) StartMarketStatusFeed() {

	storedHash := ""

	for {
		// Api call to Trader to get status
		status, err := CheckMarketStatus()
		services.Warning(err)

		// Take md5 of the status
		hash, err := structhash.Hash(status, 1)
		services.Warning(err)

		// If the hashes do not match we know the market status has changed
		if hash != storedHash {

			// Build json to send
			json, err := t.WsSendJsonBuild("change-detected", `{ "type": "market-status" }`)
			services.Warning(err)

			// Send status to all connections
			t.WsDispatchToAll(json)

			// Log event
			services.Info("StartMarketStatusFeed() : Market status has changed to " + status.State)

			// Just with this special case do we not go through the notify package. If storedHash is empty
			// we know the app just started so most likely we do not want to push a notification.
			if storedHash != "" {
				go web_push.Push(t.DB, 0, "market-status", `{ "status": "`+status.State+`"}`)
				go sms_push.Push(t.DB, 0, "market-status", `{ "status": "`+status.State+`"}`)
			}

		}

		// Store hash
		storedHash = hash

		// Save the market status in our cache.
		cache.Set("oc-market-status", status)

		// Sleep for 2 second.
		time.Sleep(time.Second * 2)
	}

}

//
// Check market status
//
func CheckMarketStatus() (MarketStatus, error) {
	var status MarketStatus

	// Setup http client
	client := &http.Client{}

	// Setup api request
	req, _ := http.NewRequest("GET", "https://api.tradier.com/v1/markets/clock", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")))

	res, err := client.Do(req)

	if err != nil {
		return status, err
	}

	// Close Body
	defer res.Body.Close()

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return status, errors.New(fmt.Sprint("CheckMarketStatus API did not return 200, It returned ", res.StatusCode))
	}

	// Read the data we got.
	body, _ := ioutil.ReadAll(res.Body)

	// Bust open the watchlist.
	var ws map[string]MarketStatus

	if err := json.Unmarshal(body, &ws); err != nil {
		return status, err
	}

	// Set the status we return.
	status = ws["clock"]

	// Return happy
	return status, nil
}

/* End File */
