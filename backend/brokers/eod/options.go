//
// Date: 2018-10-30
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-30
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package eod

import (
	"strings"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
)

//
// Get options expiration by Symbol
//
func (t *Api) GetOptionsExpirationsBySymbol(symbol string) ([]string, error) {

	expires := []string{}
	symb := strings.ToUpper(symbol)
	tmpExpires := make(map[string]bool)

	// Get a list of all options
	options, err := t.GetOptionsBySymbol(symb)

	if err != nil {
		return expires, err
	}

	// Loop through and build chain
	for _, row := range options {
		tmpExpires[row.ExpirationDate.Format("2006-01-02")] = true
	}

	// Loop through and create an array of string
	for key := range tmpExpires {
		expires = append(expires, key)
	}

	return expires, nil
}

//
// Get an options chain by expiration.
//
func (t *Api) GetOptionsChainByExpiration(symbol string, expireStr string) (types.OptionsChain, error) {

	symb := strings.ToUpper(symbol)
	expireDate := types.Date{helpers.ParseDateNoError(expireStr).UTC()}

	// New chain to return
	chain := types.OptionsChain{
		Underlying:     symb,
		ExpirationDate: expireDate,
		Puts:           []types.OptionsChainItem{},
		Calls:          []types.OptionsChainItem{},
	}

	// Get a list of all options
	options, err := t.GetOptionsBySymbol(symb)

	if err != nil {
		return chain, err
	}

	// Loop through and build chain
	for _, row := range options {

		// We only want the expire date we passed in.
		if row.ExpirationDate != expireDate {
			continue
		}

		// Append Item
		if row.OptionType == "Call" {
			chain.Calls = append(chain.Calls, row)
		} else if row.OptionType == "Put" {
			chain.Puts = append(chain.Puts, row)
		}

	}

	// Return Chain
	return chain, nil
}

/* End File */
