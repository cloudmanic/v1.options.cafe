//
// Date: 11/7/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"app.options.cafe/backend/brokers/types"
	"github.com/tidwall/gjson"
)

//
// Search for symbols or companies
//
func (t *Api) SearchBySymbolOrCompanyName(query string) ([]types.Symbol, error) {

	m := make(map[string]types.Symbol)

	// Search by Symbol
	symbs, err := t.SearchBySymbolName(query)

	if err != nil {
		return []types.Symbol{}, err
	}

	// Search by company name
	companies, err := t.SearchByCompanyName(query)

	if err != nil {
		return []types.Symbol{}, err
	}

	// Put results of both searches into a map.
	for _, row := range symbs {
		m[row.Name] = row
	}

	for _, row := range companies {
		m[row.Name] = row
	}

	// Order the companies
	mk := make([]string, len(m))
	i := 0
	for key, _ := range m {
		mk[i] = key
		i++
	}

	sort.Strings(mk)

	// Put map back into []types.Symbol format.
	var symbols []types.Symbol

	for _, row := range mk {
		symbols = append(symbols, m[row])
	}

	// Return happy
	return symbols, nil
}

//
// Search for Symbol
//
func (t *Api) SearchBySymbolName(query string) ([]types.Symbol, error) {

	// Setup http client
	client := &http.Client{}

	// Setup api request
	req, _ := http.NewRequest("GET", apiBaseUrl+"/markets/lookup?q="+query, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))

	res, err := client.Do(req)

	if err != nil {
		return []types.Symbol{}, err
	}

	// Close Body
	defer res.Body.Close()

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return []types.Symbol{}, errors.New(fmt.Sprint("Search for symbols by name API did not return 200, It returned ", res.StatusCode))
	}

	// Read the data we got.
	body, _ := ioutil.ReadAll(res.Body)

	// Parse and return.
	return parseSearchJsonResponse(string(body))
}

//
// Search for company
//
func (t *Api) SearchByCompanyName(query string) ([]types.Symbol, error) {

	// Setup http client
	client := &http.Client{}

	// Setup api request
	req, _ := http.NewRequest("GET", apiBaseUrl+"/markets/search?q="+query, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))

	res, err := client.Do(req)

	if err != nil {
		return []types.Symbol{}, err
	}

	// Close Body
	defer res.Body.Close()

	// Make sure the api responded with a 200
	if res.StatusCode != 200 {
		return []types.Symbol{}, errors.New(fmt.Sprint("Search for company by name API did not return 200, It returned ", res.StatusCode))
	}

	// Read the data we got.
	body, _ := ioutil.ReadAll(res.Body)

	// Parse and return.
	return parseSearchJsonResponse(string(body))
}

//
// Parse symbol json response.
//
func parseSearchJsonResponse(body string) ([]types.Symbol, error) {

	// Make sure we have at least one result.
	vo := gjson.Get(body, "securities.security")

	if !vo.Exists() {
		return []types.Symbol{}, nil
	}

	// Do we have more than one result
	vo = gjson.Get(body, "securities.security.symbol")

	// Setup returns array.
	var symbols []types.Symbol

	// More than one order??
	if !vo.Exists() {

		// Loop through the results.
		vo = gjson.Get(body, "securities.security")

		vo.ForEach(func(key, value gjson.Result) bool {

			// Add symbol to our array.
			symbols = append(symbols, types.Symbol{
				Name:        gjson.Get(value.String(), "symbol").String(),
				Description: gjson.Get(value.String(), "description").String(),
			})

			return true // keep iterating
		})

	} else {

		// Get JSON
		vo = gjson.Get(body, "securities.security")

		// Add symbol to our array.
		symbols = append(symbols, types.Symbol{
			Name:        gjson.Get(vo.String(), "symbol").String(),
			Description: gjson.Get(vo.String(), "description").String(),
		})

	}

	// Return happy
	return symbols, nil
}

/* End File */
