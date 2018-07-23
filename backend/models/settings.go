//
// Date: 7/22/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

type Settings struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	UserId    uint      `sql:"not null;index:UserId" json:"user_id"`

	// Put credit Spread
	StrategyPcsClosePrice float64 `json:"strategy_pcs_close_price" sql:"not null;default:0.03"`
	StrategyPcsOpenPrice  string  `json:"strategy_pcs_open_price" sql:"not null;default:'mid-point'"`
	StrategyPcsLots       uint    `json:"strategy_pcs_lots" sql:"not null;default:10"`

	// Call credit Spread
	StrategyCcsClosePrice float64 `json:"strategy_ccs_close_price" sql:"not null;default:0.03"`
	StrategyCcsOpenPrice  string  `json:"strategy_ccs_open_price" sql:"not null;default:'mid-point'"`
	StrategyCcsLots       uint    `json:"strategy_ccs_lots" sql:"not null;default:10"`

	// Put debit Spread
	StrategyPdsClosePrice float64 `json:"strategy_pds_close_price" sql:"not null;default:0.03"`
	StrategyPdsOpenPrice  string  `json:"strategy_pds_open_price" sql:"not null;default:'mid-point'"`
	StrategyPdsLots       uint    `json:"strategy_pds_lots" sql:"not null;default:10"`

	// Call debit Spread
	StrategyCdsClosePrice float64 `json:"strategy_cds_close_price" sql:"not null;default:0.03"`
	StrategyCdsOpenPrice  string  `json:"strategy_cds_open_price" sql:"not null;default:'mid-point'"`
	StrategyCdsLots       uint    `json:"strategy_cds_lots" sql:"not null;default:10"`
}

//
// Get or create by users id.
//
func (t *DB) SettingsGetOrCreateByUserId(userId uint) Settings {

	var s Settings

	// First make sure we don't already have this settings
	if t.Where("user_id = ?", userId).First(&s).RecordNotFound() {

		// Create entry.
		s = Settings{UserId: userId}

		t.Create(&s)

	}

	// Return the settings.
	return s
}

/* End File */
