//
// Date: 2/26/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package market

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/notify"
	"github.com/cloudmanic/app.options.cafe/backend/library/queue"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/library/worker"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cnf/structhash"
	"github.com/davecgh/go-spew/spew"
)

var storedHash string = ""

type MarketStatus struct {
	Date        string `json:"date"`
	State       string `json:"state"`
	Description string `json:"description"`
	NextState   string `json:"next_state"`
}

//
// Get market status.
//
func GetMarketStatus(job worker.JobRequest) error {
	var status MarketStatus

	// Setup http client
	client := &http.Client{}

	// Setup api request
	req, _ := http.NewRequest("GET", "https://api.tradier.com/v1/markets/clock", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")))

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	// Close Body
	defer res.Body.Close()

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprint("CheckMarketStatus API did not return 200, It returned ", res.StatusCode))
	}

	// Read the data we got.
	body, _ := ioutil.ReadAll(res.Body)

	// Bust open the watchlist.
	var ws map[string]MarketStatus

	if err := json.Unmarshal(body, &ws); err != nil {
		return err
	}

	// Set the status we return.
	status = ws["clock"]

	// Send message to websocket
	queue.Write("oc-websocket-write", `{"uri":"market-status","user_id":0,"body":`+helpers.JsonEncode(status)+`}`)

	// Save the market status in our cache.
	cache.Set("oc-market-status", status)

	// See if the market has changed
	DetectChange(job.DB, status)

	// Return happy
	return nil
}

//
// Detect Change
//
func DetectChange(db models.Datastore, status MarketStatus) {
	// Take md5 of the status
	hash, err := structhash.Hash(status, 1)

	if err != nil {
		services.Info(err)
	}

	// If the hashes do not match we know the market status has changed
	if hash != storedHash {

		// Log event
		services.InfoMsg("StartMarketStatusFeed() : Market status has changed to " + status.State)

		// Just with this special case do we not go through the notify package.
		s := status.State

		if status.State == "postmarket" {
			s = "closed"
		}

		// Some times s is empty.
		if (status.NextState != "premarket") && ((s == "closed") || (s == "open")) {

			spew.Dump(status)

			tDate := helpers.ParseDateNoError(status.Date)
			msg := tDate.Format("1/2/2006") + " - The market is now " + s + "."
			notify.Push(db, notify.NotifyRequest{Uri: "market-status-" + s, ShortMsg: msg, UserId: 0, Date: tDate})
		}
	}

	// Store hash
	storedHash = hash

}

/* End File */
