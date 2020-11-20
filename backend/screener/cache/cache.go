//
// Date: 2020-05-06
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//
// Used for in memory caching during screener actions.
//

package cache

import "app.options.cafe/models"

// Cache struct
type Cache struct {
	db            models.Datastore
	cachedSymbols map[string]models.Symbol
}

//
// New cache instance
//
func New(db models.Datastore) Cache {
	n := Cache{db: db}

	// Build cache of all symbols in the system
	s := db.GetAllSymbols()
	n.cachedSymbols = make(map[string]models.Symbol)

	// Loop through and build hash table
	for _, row := range s {
		n.cachedSymbols[row.ShortName] = row
	}

	// Return instance
	return n
}

/* End File */
