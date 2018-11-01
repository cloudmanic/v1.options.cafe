//
// Date: 2018-11-01
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-01
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package eod

import (
	"strings"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
)

//
// Return a quote by symbol.
//
func (t *Api) GetQuotes(symbols []string) ([]types.Quote, error) {

	quotes := []types.Quote{}

	for _, row := range symbols {

		symb := strings.ToUpper(row)

		// Get a list of all options
		_, lastQuote, err := t.GetOptionsBySymbol(symb)

		if err != nil {
			return quotes, err
		}

		// Build quote object
		quotes = append(quotes, types.Quote{
			Symbol: symb,
			Last:   lastQuote,
		})

	}

	return quotes, nil
}

/* End File */
