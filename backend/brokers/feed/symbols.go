//
// Date: 2/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package feed

import "sort"

//
// Return active symbols. This is handy because we
// sort and filter before returning.
//
func (t *Base) GetActiveSymbols() []string {
	var activeSymbols []string

	// Symbols we always want
	activeSymbols = append(activeSymbols, "$DJI")
	activeSymbols = append(activeSymbols, "SPX")
	activeSymbols = append(activeSymbols, "COMP")
	activeSymbols = append(activeSymbols, "VIX")

	// Get the watch lists for this user.
	watchList, err := t.DB.GetWatchlistsByUserId(t.User.Id)

	// Loop through the watchlists and add the symbols
	if err == nil {

		for _, row := range watchList {

			for _, row2 := range row.Symbols {

				activeSymbols = append(activeSymbols, row2.Symbol.ShortName)

			}

		}

	}

	// Add in the orders we want.
	t.muOrders.Lock()

	for _, row := range t.Orders {

		// Stock symbol
		activeSymbols = append(activeSymbols, row.Symbol)

		// Multi leg trade
		if row.NumLegs > 0 {

			for _, row2 := range row.Legs {

				activeSymbols = append(activeSymbols, row2.OptionSymbol)

			}

		}

		// Single option order
		if len(row.OptionSymbol) > 0 {
			activeSymbols = append(activeSymbols, row.OptionSymbol)
		}

	}

	t.muOrders.Unlock()

	// Clean up the list.
	activeSymbols = t.ToUpperStrings(activeSymbols)
	activeSymbols = t.RemoveDupsStrings(activeSymbols)

	// Sort the list.
	sort.Strings(activeSymbols)

	// Return the cleaned up list.
	return activeSymbols
}

/* End File */
