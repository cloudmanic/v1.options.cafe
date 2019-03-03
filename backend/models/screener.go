//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Screener struct
type Screener struct {
	Id         uint           `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	UserId     uint           `sql:"not null;index:UserId" json:"user_id"`
	BacktestId uint           `sql:"not null;index:BacktestId" json:"backtest_id"`
	Name       string         `sql:"not null" json:"name"`
	Strategy   string         `sql:"not null" json:"strategy"`
	Symbol     string         `sql:"not null" json:"symbol"`
	Items      []ScreenerItem `sql:"not null" json:"items"`
}

//
// Validate for this model.
//
func (a Screener) Validate(db Datastore, userId uint) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required.Error("The name field is required.")),
		validation.Field(&a.Strategy, validation.Required.Error("The strategy field is required.")),
		validation.Field(&a.Symbol, validation.Required.Error("The symbol field is required.")),
	)
}

//
// Get a Screeners by user id.
//
func (t *DB) GetScreenersByUserId(userId uint) ([]Screener, error) {

	var u []Screener

	if t.Where("user_id = ?", userId).Find(&u).RecordNotFound() {
		return u, errors.New("[Models:GetScreenersByUserId] Records not found (#001).")
	}

	if len(u) <= 0 {
		return u, errors.New("[Models:GetScreenersByUserId] Records not found (#002).")
	}

	// Loop through the screener and add the items
	for key := range u {
		t.Model(u[key]).Related(&u[key].Items)
	}

	// Return the Screeners.
	return u, nil
}

//
// Get a Screener by id and user id
//
func (t *DB) GetScreenerByIdAndUserId(id uint, userId uint) (Screener, error) {

	var u Screener

	if t.Where("user_id = ? AND id = ?", userId, id).First(&u).RecordNotFound() {
		return u, errors.New("[Models:GetScreenerByIdAndUserId] Record not found")
	}

	// Add in Item Lookup
	t.Model(u).Related(&u.Items) // Add in Items

	// Return the Screener.
	return u, nil
}

/* End File */
