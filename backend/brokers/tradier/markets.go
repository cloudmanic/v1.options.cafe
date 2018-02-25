//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/tidwall/gjson"
)

//
// Get historical quotes
//
func (t *Api) GetHistoricalQuotes(symbol string) ([]types.HistoryQuote, error) {

	quotes := []types.HistoryQuote{}

	// Setup http client
	client := &http.Client{}

	// Setup api request
	req, _ := http.NewRequest("GET", apiBaseUrl+"/markets/history?symbol="+symbol, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))

	res, err := client.Do(req)

	if err != nil {
		return quotes, err
	}

	// Close Body
	defer res.Body.Close()

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return quotes, errors.New(fmt.Sprint("GetHistoricalQuotes API did not return 200, It returned ", res.StatusCode))
	}

	// Read the data we got.
	body, _ := ioutil.ReadAll(res.Body)

	// Reach down in the json - Because Tradier is cray cray
	vo := gjson.Get(string(body), "history.day")

	if !vo.Exists() {
		// Return happy (just not more than one quote)
		return quotes, nil
	}

	if err := json.Unmarshal([]byte(vo.String()), &quotes); err != nil {
		return quotes, err
	}

	// Return happy
	return quotes, nil

}

/* End File */
