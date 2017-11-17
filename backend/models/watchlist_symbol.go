//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"

	"github.com/app.options.cafe/backend/library/services"
)

type WatchlistSymbol struct {
	Id          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	UserId      uint      `sql:"not null;index:UserId" json:"-"`
	WatchlistId uint      `sql:"not null;index:WatchlistId" json:"-"`
	SymbolId    uint      `sql:"not null;index:SymbolId" json:"-"`
	Order       uint      `sql:"not null" json:"-"`
	Symbol      Symbol    `json:"symbol"`
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
