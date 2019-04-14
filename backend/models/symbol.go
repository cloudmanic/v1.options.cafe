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

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

type Symbol struct {
	Id               uint      `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
	ShortName        string    `sql:"not null" json:"short_name"`
	Name             string    `sql:"not null" json:"name"`
	Type             string    `sql:"not null;type:ENUM('Equity', 'Option', 'Other');default:'Equity'" json:"type"`
	OptionUnderlying string    `json:"option_underlying"`
	OptionType       string    `json:"option_type"`
	OptionExpire     Date      `gorm:"type:date" json:"option_expire"`
	OptionStrike     float64   `json:"option_strike"`
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
// Get a symbol by short name
//
func (t *DB) GetSymbolByShortName(short string) (Symbol, error) {

	var symb Symbol

	if t.Where("short_name = ?", short).First(&symb).RecordNotFound() {
		return symb, errors.New("Unknown symbol.")
	}

	// Return happy.
	return symb, nil
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
		services.Critical(errors.New(err.Error() + "[Models:CreateNewOptionSymbol] - Unable to parse option symbol. - " + short))
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

		// Is this an option
		if sType == "Option" {

			// Get the parts of the option
			parts, err := helpers.OptionParse(short)

			if err != nil {
				services.Critical(errors.New(err.Error() + "[Models:CreateNewSymbol] - Unable to parse option symbol. - " + short))
				return Symbol{}, err
			}

			// Add in options data.
			symb.OptionUnderlying = parts.Symbol
			symb.OptionType = parts.Type
			symb.OptionExpire = Date{parts.Expire}
			symb.OptionStrike = parts.Strike
		}

		t.Create(&symb)

		// Log Symbol creation.
		services.InfoMsg("[Models:CreateNewSymbol] - Created a new Symbol entry - (" + short + ") " + name)

	}

	// Return the symbol.
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
		services.Critical(errors.New(err.Error() + "[Models:UpdateSymbol] - Unable to update symbol."))
		return symb, err
	}

	// Log Symbol creation.
	services.InfoMsg("[Models:CreateNewSymbol] - Update a new Symbol entry - (" + short + ") " + name)

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
		services.Critical(errors.New(err.Error() + "[Models:SearchSymbols] - Unable to search for symbols."))
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

//
// Get options symbol by parts
//
func (t *DB) GetOptionByParts(optionUnderlying string, optionType string, optionExpire time.Time, optionStrike float64) (Symbol, error) {

	var symbol Symbol

	// Query to see if this is a real symbol that we know.
	if t.Where("option_underlying = ? AND option_type = ? AND option_expire = ? AND option_strike = ?", optionUnderlying, optionType, optionExpire.Format("2006-01-02"), optionStrike).First(&symbol).RecordNotFound() {
		return symbol, errors.New("Unknown symbol.")
	}

	// Return the symbol
	return symbol, nil
}

//
// Load an options chain into the symbols table.
//
func (t *DB) LoadSymbolsByOptionsChain(chain types.OptionsChain) error {

	// Deal with puts.
	for _, row := range chain.Puts {
		t.CreateNewOptionSymbol(row.Symbol)
	}

	// Deal with calls.
	for _, row := range chain.Calls {
		t.CreateNewOptionSymbol(row.Symbol)
	}

	return nil
}

/* End File */
