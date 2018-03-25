//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"strings"
	"time"
)

type Watchlist struct {
	Id        uint              `gorm:"primary_key" json:"id"`
	CreatedAt time.Time         `json:"-"`
	UpdatedAt time.Time         `json:"-"`
	UserId    uint              `sql:"not null;index:UserId" json:"user_id"`
	Name      string            `sql:"not null" json:"name"`
	Symbols   []WatchlistSymbol `json:"symbols"`
}

//
// Update the name of a watchlist.
//
func (t *DB) WatchlistUpdate(id uint, name string) error {

	// Get the current watchlist
	wList, err := t.GetWatchlistsById(id)

	if err != nil {
		return err
	}

	// Update.
	wList.Name = strings.Trim(name, " ")
	t.Save(&wList)

	return nil
}

//
// Delete a watchlist by id.
//
func (t *DB) WatchlistDeleteById(id uint) error {

	// First we have to delete all the symbols connected with the watchlist.
	if err := t.Where("watchlist_id = ?", id).Delete(&WatchlistSymbol{}).Error; err != nil {
		return err
	}

	// Now lets delete this watchlist.
	if err := t.Where("id = ?", id).Delete(&Watchlist{}).Error; err != nil {
		return err
	}

	// Return happy
	return nil
}

//
// Get a Watchlists by id.
//
func (t *DB) GetWatchlistsById(id uint) (Watchlist, error) {

	var u Watchlist

	if t.Find(&u, id).RecordNotFound() {
		return u, errors.New("[Models:GetWatchlistsById] Record not found")
	}

	// Add in Symbols Lookup
	t.Model(u).Order("`order` asc").Related(&u.Symbols) // Add in Symbols

	// Add in symbols
	for key := range u.Symbols {
		t.Model(u.Symbols[key]).Related(&u.Symbols[key].Symbol)
		t.CreateActiveSymbol(u.UserId, u.Symbols[key].Symbol.ShortName)
	}

	// Return the Watchlists.
	return u, nil
}

//
// Get the watchlist by user and id. This is useful to make sure we do not give the
// watchlist to the wrong user (security).
//
func (t *DB) GetWatchlistsByIdAndUserId(id uint, userId uint) (Watchlist, error) {

	var u Watchlist

	if t.Where("user_id = ? AND id = ?", userId, id).First(&u).RecordNotFound() {
		return u, errors.New("[Models:GetWatchlistsByIdAndUserId] Record not found")
	}

	// Add in Symbols Lookup
	t.Model(u).Order("`order` asc").Related(&u.Symbols) // Add in Symbols

	// Add in symbols
	for key := range u.Symbols {
		t.Model(u.Symbols[key]).Related(&u.Symbols[key].Symbol)
		t.CreateActiveSymbol(u.UserId, u.Symbols[key].Symbol.ShortName)
	}

	// Return the Watchlists.
	return u, nil
}

//
// Get a Watchlists by user id.
//
func (t *DB) GetWatchlistsByUserId(userId uint) ([]Watchlist, error) {

	var u []Watchlist

	if t.Where("user_id = ?", userId).Find(&u).RecordNotFound() {
		return u, errors.New("[Models:GetWatchlistsByUserId] Records not found (#001).")
	}

	if len(u) <= 0 {
		return u, errors.New("[Models:GetWatchlistsByUserId] Records not found (#002).")
	}

	// Loop through the watchlist and add the items
	for key := range u {

		// Add in Symbols Lookup
		t.Model(u[key]).Order("`order` asc").Related(&u[key].Symbols)

		// Add in Symbols
		for key2 := range u[key].Symbols {
			t.Model(u[key].Symbols[key2]).Related(&u[key].Symbols[key2].Symbol)
			t.CreateActiveSymbol(u[key].UserId, u[key].Symbols[key2].Symbol.ShortName)
		}

	}

	// Return the Watchlists.
	return u, nil
}

//
// Create Watchlist entry.
//
func (t *DB) CreateWatchlist(userId uint, name string) (Watchlist, error) {

	// Create entry.
	wList := Watchlist{Name: name, UserId: userId}

	t.Create(&wList)

	// Return the user.
	return wList, nil
}

//
// Reorder a watchlist
//
func (t *DB) WatchlistReorder(id uint, ids []int) error {

	// Loop through and change the order.
	for key, row := range ids {
		t.Model(&WatchlistSymbol{}).Where("id = ? AND watchlist_id = ?", row, id).Update("order", key)
	}

	return nil
}

//
// Delete symbol from watchlist.
//
func (t *DB) WatchlistRemoveSymbol(id uint, symbId uint) error {

	t.Where("id = ? AND watchlist_id = ?", symbId, id).Delete(&WatchlistSymbol{})

	return nil
}

/* End File */
