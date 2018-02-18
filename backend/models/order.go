package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Order struct {
	Id                uint `gorm:"primary_key"`
	UserId            uint `sql:"not null;index:UserId"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	BrokerAccountId   uint   `sql:"not null;index:BrokerAccountId"`
	BrokerRef         string `sql:"not null;index:BrokerRef"`
	BrokerAccountRef  string `sql:"not null;index:BrokerAccountRef"`
	Type              string
	SymbolId          uint
	OptionSymbolId    uint
	Side              string
	Qty               int
	Status            string
	Duration          string
	Price             float64 `sql:"type:DECIMAL(12,2)"`
	AvgFillPrice      float64 `sql:"type:DECIMAL(12,2)"`
	ExecQuantity      float64 `sql:"type:DECIMAL(12,2)"`
	LastFillPrice     float64 `sql:"type:DECIMAL(12,2)"`
	LastFillQuantity  float64 `sql:"type:DECIMAL(12,2)"`
	RemainingQuantity float64 `sql:"type:DECIMAL(12,2)"`
	CreateDate        time.Time
	TransactionDate   time.Time
	Class             string
	NumLegs           int
	PositionReviewed  string `sql:"not null;type:ENUM('No', 'Yes', 'Error');default:'No'"`
	Symbol            Symbol `json:"symbol"`
	OptionSymbol      Symbol `gorm:"foreignkey:OptionSymbolId" json:"option_symbol"`
	Legs              []OrderLeg
}

type OrderLeg struct {
	Id                uint `gorm:"primary_key"`
	UserId            uint `sql:"not null;index:UserId"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	OrderId           uint `sql:"not null;index:OrderId"`
	Type              string
	SymbolId          uint
	Side              string
	Qty               int
	Status            string
	Duration          string
	AvgFillPrice      float64 `sql:"type:DECIMAL(12,2)"`
	ExecQuantity      float64 `sql:"type:DECIMAL(12,2)"`
	LastFillPrice     float64 `sql:"type:DECIMAL(12,2)"`
	LastFillQuantity  float64 `sql:"type:DECIMAL(12,2)"`
	RemainingQuantity float64 `sql:"type:DECIMAL(12,2)"`
	CreateDate        time.Time
	TransactionDate   time.Time
	Symbol            Symbol `json:"symbol"`
}

//
// Store a new order.
//
func (t *DB) CreateOrder(order *Order) error {

	// Create order
	t.Create(order)

	// Return happy
	return nil
}

//
// Store a order.
//
func (t *DB) UpdateOrder(order *Order) error {

	// Update entry.
	t.Save(&order)

	// Return happy
	return nil
}

//
// Get orders by User and class and status and reviewed
//
func (t *DB) GetOrdersByUserClassStatusReviewed(userId uint, class string, status string, reviewed string) ([]Order, error) {

	orders := []Order{}

	// Query and get all orders we have not reviewed before.
	t.Debug().Preload("Legs", func(db *gorm.DB) *gorm.DB {
		return db.Order("id asc")
	}).Where("user_id = ? AND class = ? AND status = ? AND position_reviewed = ?", userId, class, status, reviewed).Order("transaction_date asc").Find(&orders)

	// Return happy
	return orders, nil
}

//
// See if we have an order by user and broker ref.
//
func (t *DB) HasOrderByBrokerRefUserId(brokerId string, userId uint) bool {

	// See if we already have this record in our database
	var count int
	order := &Order{}

	// Run query
	t.Where("broker_ref = ? AND user_id = ?", brokerId, userId).First(order).Count(&count)

	if count > 0 {
		return true
	}

	// Return not found.
	return false
}

/* End File */
