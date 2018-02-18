package models

import (
	"errors"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

type TradeGroup struct {
	Id               uint       `gorm:"primary_key" json:"id"`
	UserId           uint       `sql:"not null;index:UserId" json:"_"`
	CreatedAt        time.Time  `json:"_"`
	UpdatedAt        time.Time  `json:"_"`
	Name             string     `json:"name"`
	BrokerAccountId  uint       `sql:"not null;index:BrokerAccountId" json:"_"`
	BrokerAccountRef string     `sql:"not null;index:BrokerAccountRef" json:"broker_account_ref"`
	Status           string     `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'" json:"status"`
	Type             string     `sql:"not null;type:ENUM('Option', 'Stock', 'Put Credit Spread', 'Call Credit Spread', 'Put Debit Spread', 'Call Debit Spread', 'Iron Condor', 'Other'); default:'Other'" json:"type"`
	OrderIds         string     `json:"_"`
	Risked           float64    `sql:"type:DECIMAL(12,2)" json:"risked"`
	Credit           float64    `sql:"type:DECIMAL(12,2)" json:"credit"`       // Before Commission
	Proceeds         float64    `sql:"type:DECIMAL(12,2)" json:"proceeds"`     // Before Commission
	Profit           float64    `sql:"type:DECIMAL(12,2)" json:"profit"`       // After Commission
	PercentGain      float64    `sql:"type:DECIMAL(12,2)" json:"percent_gain"` // After Commission
	Commission       float64    `sql:"type:DECIMAL(12,2)" json:"commission"`
	Note             string     `sql:"type:text" json:"note"`
	Positions        []Position `json:"positions"`
	OpenDate         time.Time  `json:"open_date"`
	ClosedDate       time.Time  `json:"closed_date"`
}

//
// List different trade groups.
//
func (t *DB) GetTradeGroups(params QueryParam) ([]TradeGroup, QueryMetaData, error) {

	var results = []TradeGroup{}

	// Run the query
	noFilterCount, err := t.QueryWithNoFilterCount(&results, params)

	// Throw error if we have one
	if err != nil {
		return results, QueryMetaData{}, err
	}

	// Get the meta data related to this query.
	meta := t.GetQueryMetaData(len(results), noFilterCount, params)

	// Add the symbol data to the positions.
	t.tradeGroupAddSymbolsToPositions(results)

	// Return happy
	return results, meta, nil
}

//
// Get TradeGroup by Id
//
func (t *DB) GetTradeGroupById(id uint) (TradeGroup, error) {

	tg := TradeGroup{}

	if t.Preload("Positions").Where("Id = ?", id).First(&tg).RecordNotFound() {
		return tg, errors.New("Record not found")
	}

	// Loop through and add the symbol to the positions object
	for key, row := range tg.Positions {
		t.Model(row).Related(&tg.Positions[key].Symbol)
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

//
// Helpful function for adding symbols to positions
//
func (t *DB) tradeGroupAddSymbolsToPositions(tgs []TradeGroup) error {

	// Loop through and add the symbol to the positions object
	for key, row := range tgs {
		for key2, row2 := range tgs[key].Positions {
			t.Model(row2).Related(&tgs[key].Positions[key2].Symbol)
		}

		// Stupid code to make it so GoFMT does not delete in range statement
		if row.Id < 0 {
			services.Info("This should never happen.")
		}
	}

	return nil
}

/* End File */
