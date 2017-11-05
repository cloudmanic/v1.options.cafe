//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"

	"app.options.cafe/backend/library/services"
)

type WatchlistSymbol struct {
	Id          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserId      uint `sql:"not null;index:UserId"`
	WatchlistId uint `sql:"not null;index:WatchlistId"`
	SymbolId    uint `sql:"not null;index:SymbolId"`
	Order       uint `sql:"not null"`
	Symbol      Symbol
}

//
// Create a new WatchlistSymbol entry.
//
func (t *DB) CreateNewWatchlistSymbol(wList Watchlist, symb Symbol, user User, order uint) (WatchlistSymbol, error) {

	// Create entry.
	wListSym := WatchlistSymbol{WatchlistId: wList.Id, SymbolId: symb.Id, UserId: user.Id, Order: order}

	t.Create(&wListSym)

	// Log broker creation.
	services.Log("CreateNewWatchlistSymbol - Created a new WatchlistSymbol entry - " + symb.ShortName + " " + wList.Name)

	// Return the user.
	return wListSym, nil

}

/* End File */
