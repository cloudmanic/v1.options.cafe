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

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

type Symbol struct {
	Id            uint          `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time     `json:"-"`
	UpdatedAt     time.Time     `json:"-"`
	ShortName     string        `sql:"not null" json:"short_name"`
	Name          string        `sql:"not null" json:"name"`
	Type          string        `sql:"not null;type:ENUM('Equity', 'Option', 'Other');default:'Equity'" json:"type"`
	OptionDetails OptionDetails `json:"option_details"`
}

type OptionDetails struct {
	Symbol string    `sql:"not null" json:"symbol"`
	Expire time.Time `json:"expire"`
	Type   string    `json:"type"`
	Strike float64   `json:"strike"`
}

//
// Validate if a symbol is real.
//
func (t *DB) ValidateSymbolId(value interface{}) error {

	// Query to see if this is a real symbol that we know.
	if t.Where("Id = ?", value).First(&Symbol{}).RecordNotFound() {
		return errors.New("Unknown symbol_id.")
	}

	return nil
}

//
// Special function just to create a new options symbol.
// We pass in a option string such as "SPY180209P00276000"
// We then build the full name of the option.
//
func (t *DB) CreateNewOptionSymbol(short string) (Symbol, error) {

	// Get the parts of the option
	parts, err := helpers.OptionParse(short)

	if err != nil {
		services.Error(err, "[Models:CreateNewOptionSymbol] - Unable to parse option symbol.")
		return Symbol{}, err
	}

	// Store the symbol
	return t.CreateNewSymbol(short, parts.Name, "Option")
}

//
// Create a new Symbol entry.
//
func (t *DB) CreateNewSymbol(short string, name string, sType string) (Symbol, error) {

	var symb Symbol

	// First make sure we don't already have this symbol
	if t.Where("short_name = ?", short).First(&symb).RecordNotFound() {

		// Create entry.
		symb = Symbol{Name: name, ShortName: strings.ToUpper(short), Type: sType}

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
func (t *DB) UpdateSymbol(id uint, short string, name string, sType string) (Symbol, error) {

	var symb Symbol

	t.First(&symb, id)
	symb.Name = name
	symb.ShortName = strings.ToUpper(short)
	symb.Type = sType
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

	// Make query.
	t.Find(&symbols)

	// Add in options parts
	for key, row := range symbols {
		if row.Type == "Option" {
			parts, _ := helpers.OptionParse(row.ShortName)
			symbols[key].OptionDetails.Type = parts.Type
			symbols[key].OptionDetails.Strike = parts.Strike
			symbols[key].OptionDetails.Expire = parts.Expire
		}
	}

	return symbols
}

//
// Search for symbols by query string.
//
func (t *DB) SearchSymbols(query string, sType string) ([]Symbol, error) {

	var symbols []Symbol

	var sql = `SELECT *,
    IF(short_name = ?,  40, IF(short_name LIKE ?, 10, 0))
	  + IF(short_name LIKE ?, 20,  0)    
		+ IF(name LIKE ?, 10,  0)    
    + IF(name LIKE ?, 5,  0)
    AS weight
		FROM symbols
		WHERE (short_name LIKE ? OR name LIKE ?) AND type=?
		ORDER BY weight DESC
		LIMIT 10`

	rows, err := t.Raw(sql, query, "%"+query+"%", query+"%", query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", sType).Rows()

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
