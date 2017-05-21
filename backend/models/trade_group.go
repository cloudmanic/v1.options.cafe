package models

import (
  "time"
)

type TradeGroup struct {
  Id uint `gorm:"primary_key"`
  UserId uint `sql:"not null;index:UserId"`
  CreatedAt time.Time
  UpdatedAt time.Time
  AccountId string `sql:"not null;index:AccountId"`
  Status string `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'"`
  OrderIds string
  Risked float64 `sql:"type:DECIMAL(12,2)"`
  Gain float64 `sql:"type:DECIMAL(12,2)"` // Before Commission
  Profit float64 `sql:"type:DECIMAL(12,2)"` // After Commission
  Commission float64 `sql:"type:DECIMAL(12,2)"`
  Note string `sql:"type:text"`
  OpenDate time.Time
  ClosedDate time.Time
}

/* End File */