//
// Date: 2/26/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

type ActiveSymbol struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	UserId    uint      `sql:"not null;index:UserId" json:"-"`
	Symbol    string    `sql:"not null;" json:"symbol"`
}

//
// Get active symbols by user
//
func (t *DB) GetActiveSymbolsByUser(userId uint) ([]ActiveSymbol, error) {

	var result = []ActiveSymbol{}

	t.Where(ActiveSymbol{UserId: userId}).Find(&result)

	return result, nil
}

//
// Create a new Active Symbol entry.
//
func (t *DB) CreateActiveSymbol(userId uint, symbol string) (ActiveSymbol, error) {

	entry := ActiveSymbol{}

	t.Where(ActiveSymbol{UserId: userId, Symbol: symbol}).FirstOrCreate(&entry)

	return entry, nil
}

/* End File */
