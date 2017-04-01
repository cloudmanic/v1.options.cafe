package models

import (
  "time"
)

type Position struct {
  Id uint `gorm:"primary_key"`
  UserId uint `sql:"not null;index:UserId"`
  TradeGroupId uint `sql:"not null;index:TradeGroupId"`
  CreatedAt time.Time
  UpdatedAt time.Time
  AccountId string `sql:"not null;index:AccountId"`
  Status string `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'"`
  Symbol string
  Qty int
  OrgQty int
  CostBasis float64 `sql:"type:DECIMAL(12,2)"`
  AvgOpenPrice float64 `sql:"type:DECIMAL(12,2)"`
  AvgClosePrice float64 `sql:"type:DECIMAL(12,2)"`
  OrderIds string
  Note string `sql:"type:text"`
  OpenDate time.Time
  ClosedDate time.Time
}

/* End File */