//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

//
// Get time sales data
//
// TODO: pull the access token from the user's account.
//
func (t *Controller) GetHistoricalQuotes(c *gin.Context) {

	// Setup http client
	client := &http.Client{}

	// Setup api request
	req, _ := http.NewRequest("GET", "https://api.tradier.com/v1/markets/history?symbol=spy", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")))

	res, err := client.Do(req)

	// Throw error if we have one
	if t.RespondError(c, err, httpGenericErrMsg) {
		return
	}
	// Close Body
	defer res.Body.Close()

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		t.RespondError(c, errors.New(fmt.Sprint("/markets/history API did not return 200, It returned ", res.StatusCode)), httpGenericErrMsg)
		return
	}

	// Read the data we got.
	body, _ := ioutil.ReadAll(res.Body)

	spew.Dump(string(body))

	// // Bust open the watchlist.
	// var ws map[string]types.MarketStatus

	// if err := json.Unmarshal(body, &ws); err != nil {
	// 	return status, err
	// }

	// // Set the status we return.
	// status = ws["clock"]

	// // Return happy
	// return status, nil

	// Return happy JSON
	//c.JSON(200, res.Body)
}

/* End File */
