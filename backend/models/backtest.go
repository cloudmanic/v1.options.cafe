//
// Date: 2/22/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Backtest struct
type Backtest struct {
	Id               uint                 `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time            `json:"-"`
	UpdatedAt        time.Time            `json:"-"`
	UserId           uint                 `sql:"not null;index:UserId" json:"user_id"`
	StartDate        Date                 `gorm:"type:date" sql:"not null" json:"start_date"`
	EndDate          Date                 `gorm:"btype:date" sql:"not null" json:"end_date"`
	EndingBalance    float64              `sql:"not null" json:"ending_balance"`
	StartingBalance  float64              `sql:"not null" json:"starting_balance"`
	CAGR             float64              `sql:"not null" json:"cagr"`
	Return           float64              `sql:"not null" json:"return"`
	Profit           float64              `sql:"not null" json:"profit"`
	TradeCount       int                  `sql:"not null" json:"trade_count"`
	TradeSelect      string               `sql:"not null;type:ENUM('least-days-to-expire', 'highest-midpoint', 'highest-ask', 'highest-percent-away', 'shortest-percent-away');default:'highest-midpoint'" json:"trade_select"`
	Midpoint         bool                 `sql:"not null" json:"midpoint"` // Open trade at the midpoint
	PositionSize     string               `sql:"not null" json:"position_size"`
	TimeElapsed      time.Duration        `sql:"not null" json:"time_elapsed"`
	Benchmark        string               `sql:"not null" json:"benchmark"`
	BenchmarkStart   float64              `sql:"not null" json:"benchmark_start"`
	BenchmarkEnd     float64              `sql:"not null" json:"benchmark_end"`
	BenchmarkCAGR    float64              `sql:"not null" json:"benchmark_cagr"`
	BenchmarkPercent float64              `sql:"not null" json:"benchmark_percent"`
	Screen           Screener             `json:"screen"`
	TradeGroups      []BacktestTradeGroup `json:"trade_groups"`
}

//
// Validate for this model.
//
func (a Backtest) Validate(db Datastore, userId uint) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.StartingBalance, validation.Required.Error("The starting balance field is required.")),

		validation.Field(&a.TradeSelect, validation.Required.Error("The trade select field is required.")),

		validation.Field(&a.PositionSize, validation.Required.Error("The position size field is required.")),

		validation.Field(&a.Benchmark, validation.Required.Error("The benchmark field is required.")),

		validation.Field(&a.Screen,
			validation.By(func(value interface{}) error { return db.ValidateBacktestScreener(a, userId) }),
		),
	)
}

//
// ValidateBacktestScreener - Make sure all is good.
//
func (db *DB) ValidateBacktestScreener(backtest Backtest, userId uint) error {
	const errMsg1 = "Screener Id is required."

	_, err := db.GetScreenerByIdAndUserId(backtest.Screen.Id, userId)

	if err != nil {
		return errors.New(errMsg1)
	}

	// All good in the hood
	return nil
}

//
// BacktestsGetByUserId returns all the backtests for a user.
//
func (t *DB) BacktestsGetByUserId(userId uint) ([]Backtest, error) {
	bts := []Backtest{}

	if t.Where("user_id = ?", userId).Preload("Screen").Preload("Screen.Items").Preload("TradeGroups").Preload("TradeGroups.Positions").Preload("TradeGroups.Positions.Symbol").Find(&bts).RecordNotFound() {
		return bts, errors.New("[Models:BacktestsGetByUserId] Records not found (#001).")
	}

	// Return happy
	return bts, nil
}

//
// BacktestGetById returns a backtest by id.
//
func (t *DB) BacktestGetById(id uint) (Backtest, error) {
	bt := Backtest{}

	if t.Preload("Screen").Preload("Positions").Preload("Screen.Items").Preload("Positions.Legs").Where("Id = ?", id).First(&bt).RecordNotFound() {
		return bt, errors.New("Record not found")
	}

	// Return happy
	return bt, nil
}

//
// CreateBacktest will create an entry.
//
func (t *DB) CreateBacktest(backtest Backtest) (Backtest, error) {
	t.Create(&backtest)
	return backtest, nil
}

/* End File */
