//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"strings"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

type Symbol struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	ShortName string    `sql:"not null" json:"short_name"`
	Name      string    `sql:"not null" json:"name"`
}

//
// Create a new Symbol entry.
//
func (t *DB) CreateNewSymbol(short string, name string) (Symbol, error) {

	var symb Symbol

	// First make sure we don't already have this symbol
	if t.Where("short_name = ?", short).First(&symb).RecordNotFound() {

		// Create entry.
		symb = Symbol{Name: name, ShortName: strings.ToUpper(short)}

		t.Create(&symb)

		// Log Symbol creation.
		services.Info("[Models:CreateNewSymbol] - Created a new Symbol entry - (" + short + ") " + name)

	}

	// Return the user.
	return symb, nil

}

//
// Update Symbol entry.
//
func (t *DB) UpdateSymbol(id uint, short string, name string) (Symbol, error) {

	var symb Symbol

	t.First(&symb, id)
	symb.Name = name
	symb.ShortName = strings.ToUpper(short)
	err := t.Save(&symb).Error

	if err != nil {
		services.Error(err, "[Models:UpdateSymbol] - Unable to update symbol.")
		return symb, err
	}

	// Log Symbol creation.
	services.Info("[Models:CreateNewSymbol] - Update a new Symbol entry - (" + short + ") " + name)

	// Return the user.
	return symb, nil
}

//
// Get all symbols.
//
func (t *DB) GetAllSymbols() []Symbol {

	var symbols []Symbol

	t.Find(&symbols)

	return symbols
}

//
// Search for symbols by query string.
//
func (t *DB) SearchSymbols(query string) ([]Symbol, error) {

	var symbols []Symbol

	var sql = `SELECT *,
    IF(short_name = ?,  40, IF(short_name LIKE ?, 10, 0))
	  + IF(short_name LIKE ?, 20,  0)    
		+ IF(name LIKE ?, 10,  0)    
    + IF(name LIKE ?, 5,  0)
    AS weight
		FROM symbols
		WHERE (short_name LIKE ? OR name LIKE ?)
		ORDER BY weight DESC
		LIMIT 10`

	rows, err := t.Raw(sql, query, "%"+query+"%", query+"%", query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").Rows()

	if err != nil {
		services.Error(err, "[Models:SearchSymbols] - Unable to search for symbols.")
		return symbols, err
	}

	defer rows.Close()

	for rows.Next() {

		var s Symbol
		t.ScanRows(rows, &s)

		symbols = append(symbols, s)
	}

	return symbols, nil
}

/* End File */
