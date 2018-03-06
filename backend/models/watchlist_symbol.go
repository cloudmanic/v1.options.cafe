//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/davecgh/go-spew/spew"
	validation "github.com/go-ozzo/ozzo-validation"
)

type WatchlistSymbol struct {
	Id          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	UserId      uint      `sql:"not null;index:UserId" json:"-"`
	WatchlistId uint      `sql:"not null;index:WatchlistId" json:"watchlist_id"`
	SymbolId    uint      `sql:"not null;index:SymbolId" json:"symbol_id"`
	Order       uint      `sql:"not null" json:"-"`
	Symbol      Symbol    `json:"symbol"`
}

//
// Validate for this model.
//
func (a WatchlistSymbol) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.WatchlistId, validation.Required.Error("The watchlist_id field is required.")),
		validation.Field(&a.SymbolId, validation.Required.Error("The symbol_id field is required.")),
		validation.Field(&a.UserId, validation.Required.Error("The user_id field is required.")),
	)
}

//
// Prepend a new WatchlistSymbol entry.
//
func (t *DB) PrependWatchlistSymbol(wlistSymb *WatchlistSymbol) error {

	spew.Dump(wlistSymb)

	// // Create entry.
	// wListSym := WatchlistSymbol{WatchlistId: wList.Id, SymbolId: symb.Id, UserId: user.Id, Order: order}

	// t.Create(&wListSym)

	// // Store in active
	// t.CreateActiveSymbol(user.Id, symb.ShortName)

	// // Log broker creation.
	// services.Info("CreateNewWatchlistSymbol - Created a new WatchlistSymbol entry - " + symb.ShortName + " " + wList.Name)

	// Return the user.
	return nil
}

//
// Create a new WatchlistSymbol entry.
//
func (t *DB) CreateWatchlistSymbol(wList Watchlist, symb Symbol, user User, order uint) (WatchlistSymbol, error) {

	// Create entry.
	wListSym := WatchlistSymbol{WatchlistId: wList.Id, SymbolId: symb.Id, UserId: user.Id, Order: order}

	t.Create(&wListSym)

	// Store in active
	t.CreateActiveSymbol(user.Id, symb.ShortName)

	// Log broker creation.
	services.Info("CreateNewWatchlistSymbol - Created a new WatchlistSymbol entry - " + symb.ShortName + " " + wList.Name)

	// Return the user.
	return wListSym, nil
}

/* End File */
