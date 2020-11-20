//
// Date: 2020-05-06
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//
// Used for in memory caching during screener actions.
//

package cache

import "app.options.cafe/models"

//
// GetSymbol - This is a wrapper function for models.Symbol. We want to do a bit of "caching".
//
func (t *Cache) GetSymbol(short string, name string, sType string) (models.Symbol, error) {
	// See if we have this symbol in cache. If so return happy.
	if val, ok := t.cachedSymbols[short]; ok {
		return val, nil
	}

	// Add symbol to the DB. Since we do not know about it.
	symb, err := t.db.CreateNewSymbol(short, name, sType)

	if err != nil {
		return symb, err
	}

	// Add symbol to map. TODO(spicer) have to do magic to manage concurrent access with go routines
	//t.CachedSymbols[short] = symb

	// Return happy.
	return symb, nil
}

/* End File */
