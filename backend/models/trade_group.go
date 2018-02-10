package models

import (
	"errors"
	"time"
)

type TradeGroup struct {
	Id              uint       `gorm:"primary_key" json:"id"`
	UserId          uint       `sql:"not null;index:UserId" json:"_"`
	CreatedAt       time.Time  `json:"_"`
	UpdatedAt       time.Time  `json:"_"`
	BrokerAccountId uint       `sql:"not null;index:BrokerAccountId" json:"_"`
	AccountId       string     `sql:"not null;index:AccountId" json:"account_id"`
	Status          string     `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'" json:"status"`
	Type            string     `sql:"not null;type:ENUM('Option', 'Stock', 'Put Credit Spread', 'Call Credit Spread', 'Put Debit Spread', 'Call Debit Spread', 'Iron Condor', 'Other'); default:'Other'" json:"type"`
	OrderIds        string     `json:"_"`
	Risked          float64    `sql:"type:DECIMAL(12,2)" json:"risked"`
	Gain            float64    `sql:"type:DECIMAL(12,2)" json:"gain"`   // Before Commission
	Profit          float64    `sql:"type:DECIMAL(12,2)" json:"profit"` // After Commission
	Commission      float64    `sql:"type:DECIMAL(12,2)" json:"commission"`
	Note            string     `sql:"type:text" json:"note"`
	Positions       []Position `json:"positions"`
	OpenDate        time.Time  `json:"open_date"`
	ClosedDate      time.Time  `json:"closed_date"`
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
// Get a TradeGroup by user id.
//
func (t *DB) GetTradeGroupsByUserId(userId uint, orderBy string) ([]TradeGroup, error) {

	var u []TradeGroup

	if t.Preload("Positions").Where("user_id = ?", userId).Order(orderBy).Find(&u).RecordNotFound() {
		return u, errors.New("[Models:GetTradeGroupByUserId] Records not found (#001).")
	}

	if len(u) <= 0 {
		return u, errors.New("[Models:GetTradeGroupByUserId] Records not found (#002).")
	}

	// Return the TradeGroups.
	return u, nil
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
