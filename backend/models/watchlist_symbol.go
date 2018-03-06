//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
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
func (a WatchlistSymbol) Validate(db Datastore) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.WatchlistId, validation.Required.Error("The watchlist_id field is required.")),
		validation.Field(&a.SymbolId, validation.Required.Error("The symbol_id field is required."), validation.By(db.ValidateSymbolId)),
		validation.Field(&a.UserId, validation.Required.Error("The user_id field is required.")),
	)
}

//
// Prepend a new WatchlistSymbol entry.
//
func (t *DB) PrependWatchlistSymbol(w *WatchlistSymbol) error {

	// Just double check this symbol is not already in the database for this watchlist.
	if !t.Where("watchlist_id = ? AND symbol_id = ?", w.WatchlistId, w.SymbolId).Find(&WatchlistSymbol{}).RecordNotFound() {
		return errors.New("Symbol already part of this watchlist.")
	}

	// List of watchlists.
	list := []WatchlistSymbol{}

	// Loop through the watch list and move the order of each symbol down by one in order.
	t.Where("watchlist_id = ?", w.WatchlistId).Find(&list)

	for _, row := range list {
		row.Order++
		t.Save(&row)
	}

	// Insert the new symbol at the top.
	w.Order = 0
	t.Save(w)

	// We load in the symbol data for fun (as we return the full object via the API)
	t.Model(&w).Related(&w.Symbol)

	// Add in parts (even tho options in a watchlist is odd)
	parts, _ := helpers.OptionParse(w.Symbol.ShortName)
	w.Symbol.OptionDetails.Type = parts.Type
	w.Symbol.OptionDetails.Strike = parts.Strike
	w.Symbol.OptionDetails.Expire = parts.Expire

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
