//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"

	"app.options.cafe/backend/library/services"
)

type Watchlist struct {
	Id        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    uint   `sql:"not null;index:UserId"`
	Name      string `sql:"not null"`
	Symbols   []WatchlistSymbol
}

//
// Get a Watchlists by user id.
//
func (t *DB) GetWatchlistsByUserId(userId uint) ([]Watchlist, error) {

	var u []Watchlist

	if t.Where("user_id = ?", userId).Find(&u).RecordNotFound() {
		return u, errors.New("Records not found")
	}

	if len(u) <= 0 {
		return u, errors.New("Records not found")
	}

	// Loop through the watchlist and add the items
	for key := range u {

		// Add in Symbols Lookup
		t.Model(u[key]).Related(&u[key].Symbols)

		// Add in Symbols
		for key2 := range u[key].Symbols {
			t.Model(u[key].Symbols[key2]).Related(&u[key].Symbols[key2].Symbol)
		}

	}

	// Return the Watchlists.
	return u, nil

}

//
// Create a new Watchlist entry.
//
func (t *DB) CreateNewWatchlist(user User, name string) (Watchlist, error) {

	// Create entry.
	wList := Watchlist{Name: name, UserId: user.Id}

	t.Create(&wList)

	// Log broker creation.
	services.Log("CreateNewWatchlist - Created a new Watchlist entry - " + name + " " + user.Email)

	// Return the user.
	return wList, nil

}

/* End File */
