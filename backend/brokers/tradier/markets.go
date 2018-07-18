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
// Make an API call to Tradier to get an option chain by expiration.
//
func (t *Api) GetOptionsChainByExpiration(symbol string, expireDate string) (types.OptionsChain, error) {

	// Get JSON data from Tradier.
	jsonResponse, err := t.SendGetRequest("/markets/options/chains?symbol=" + symbol + "&expiration=" + expireDate)

	if err != nil {
		return types.OptionsChain{}, err
	}

	// Loop through the chain (NOTE: we assume this chain has more than one option, this could be a source of bugs)
	chainJson := gjson.Get(string(jsonResponse), "options.option").String()

	// Build an array of option chain items.
	chainItems := []types.OptionsChainItem{}

	if err := json.Unmarshal([]byte(chainJson), &chainItems); err != nil {
		return types.OptionsChain{}, err
	}

	// Build the chain
	chain := types.OptionsChain{
		Underlying:     chainItems[0].Underlying,
		ExpirationDate: chainItems[0].ExpirationDate,
	}

	for _, row := range chainItems {

		// Put or a call?
		if row.OptionType == "put" {
			chain.Puts = append(chain.Puts, row)
		} else if row.OptionType == "call" {
			chain.Calls = append(chain.Calls, row)
		}

	}

	// Return happy JSON
	return chain, nil
}

//
// Get time sale quotes
// Interval : tick, 1min, 5min or 15min (default: tick)
//
func (t *Api) GetTimeSalesQuotes(symbol string, start time.Time, end time.Time, interval string) ([]types.HistoryQuote, error) {

	quotes := []types.HistoryQuote{}

	// Setup http client
	client := &http.Client{}

	// Setup api request
	req, _ := http.NewRequest("GET", apiBaseUrl+"/markets/timesales", nil)
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

//
// Get option expirations by date (this is not part of the interface)
//
func (t *Api) GetOptionsExpirationsBySymbol(symb string) ([]string, error) {

	var result []string

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", apiBaseUrl+"/markets/options/expirations?symbol="+symb, nil)

	// Headers
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))
	req.Header.Add("Accept", "application/json")

	// Fetch Request
	res, err := client.Do(req)

	if err != nil {
		return result, err
	}

	// Close Body
	defer res.Body.Close()

	// Read Response Body
	json, _ := ioutil.ReadAll(res.Body)

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return result, errors.New("Failed response from Tradier.")
	}

	// Loop through the dates
	dates := gjson.Get(string(json), "expirations.date")
	for _, row := range dates.Array() {
		result = append(result, row.String())
	}

	// Return happy
	return result, nil
}

/* End File */
