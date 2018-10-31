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
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/tidwall/gjson"
)

//
// Get time sale quotes
// Interval : tick, 1min, 5min or 15min (default: tick)
//
func (t *Api) GetTimeSalesQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error) {

	quotes := []types.HistoryQuote{}

	// Setup http client
	client := &http.Client{}

	// Get url to api
	apiUrl := apiBaseUrl

	if t.Sandbox {
		apiUrl = sandBaseUrl
	}

	// Setup api request
	req, _ := http.NewRequest("GET", apiUrl+"/markets/timesales", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))

	// Build query string
	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("start", start.Format("2006-01-02 15:04"))
	q.Add("end", end.Format("2006-01-02 15:04"))
	q.Add("interval", interval)
	q.Add("session_filter", "open")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)

	if err != nil {
		return quotes, err
	}

	// Close Body
	defer res.Body.Close()

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return quotes, errors.New(fmt.Sprint("GetTimeSalesQuotes API did not return 200, It returned ", res.StatusCode))
	}

	// Read the data we got.
	body, _ := ioutil.ReadAll(res.Body)

	// Reach down in the json - Because Tradier is cray cray
	vo := gjson.Get(string(body), "series.data")

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

//
// Get historical quotes
//
func (t *Api) GetHistoricalQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error) {

	quotes := []types.HistoryQuote{}

	// Setup http client
	client := &http.Client{}

	// Build request
	request := apiBaseUrl + "/markets/history?symbol=" + symbol + "&start=" + start.Format("2006-01-02") + "&end=" + end.Format("2006-01-02") + "&interval=" + interval

	// Setup api request
	req, _ := http.NewRequest("GET", request, nil)
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
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
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
