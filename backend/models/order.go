package models

import (
	"time"
)

type Order struct {
	Id                uint `gorm:"primary_key"`
	UserId            uint `sql:"not null;index:UserId"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	BrokerId          int    `sql:"not null;index:OurBrokerId"`
	BrokerRef         string `sql:"not null;index:BrokerRef"`
	AccountId         string `sql:"not null;index:AccountId"`
	Type              string
	Symbol            string
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
	PositionReviewed  string `sql:"not null;type:ENUM('No', 'Yes');default:'No'"`
	Legs              []OrderLeg
}

type OrderLeg struct {
	Id                uint `gorm:"primary_key"`
	UserId            uint `sql:"not null;index:UserId"`
	OrderId           uint `sql:"not null;index:OrderId"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Type              string
	Symbol            string
	OptionSymbol      string
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
	t.Preload("Legs").Where("user_id = ? AND class = ? AND status = ? AND position_reviewed = ?", userId, class, status, reviewed).Order("transaction_date asc").Find(&orders)

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
