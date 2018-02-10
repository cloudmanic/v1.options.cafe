package models

import (
	"errors"
	"time"
)

type TradeGroup struct {
	Id              uint `gorm:"primary_key"`
	UserId          uint `sql:"not null;index:UserId"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	BrokerAccountId uint   `sql:"not null;index:BrokerAccountId"`
	AccountId       string `sql:"not null;index:AccountId"`
	Status          string `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'"`
	Type            string `sql:"not null;type:ENUM('Option', 'Stock', 'Put Credit Spread', 'Call Credit Spread', 'Put Debit Spread', 'Call Debit Spread', 'Iron Condor', 'Other'); default:'Other'"`
	OrderIds        string
	Risked          float64 `sql:"type:DECIMAL(12,2)"`
	Gain            float64 `sql:"type:DECIMAL(12,2)"` // Before Commission
	Profit          float64 `sql:"type:DECIMAL(12,2)"` // After Commission
	Commission      float64 `sql:"type:DECIMAL(12,2)"`
	Note            string  `sql:"type:text"`
	OpenDate        time.Time
	ClosedDate      time.Time
}

//
// Get TradeGroup by Id
//
func (t *DB) GetTradeGroupById(id uint) (TradeGroup, error) {

	tg := TradeGroup{}

	if t.Where("Id = ?", id).First(&tg).RecordNotFound() {
		return tg, errors.New("Record not found")
	}

	// Return happy
	return tg, nil
}

//
// Store a new TradeGroup.
//
func (t *DB) CreateTradeGroup(tg *TradeGroup) error {

	// Create order
	t.Create(tg)

	// Return happy
	return nil
}

//
// Update a TradeGroup.
//
func (t *DB) UpdateTradeGroup(tg *TradeGroup) error {

	// Update entry.
	t.Save(&tg)

	// Return happy
	return nil
}

/* End File */
