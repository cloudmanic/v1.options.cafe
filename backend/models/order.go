package models

import (
  "time"
)

type Order struct {
  Id uint `gorm:"primary_key"`
  UserId uint `sql:"not null;index:UserId"`
  CreatedAt time.Time
  UpdatedAt time.Time
  BrokerId int `sql:"not null;index:BrokerId"`
  AccountId string `sql:"not null;index:AccountId"`  
  Type string
  Symbol string
  Side string
  Qty int
  Status string
  Duration string
  Price float64 `sql:"type:DECIMAL(12,2)"`
  AvgFillPrice float64 `sql:"type:DECIMAL(12,2)"` 
  ExecQuantity float64 `sql:"type:DECIMAL(12,2)"` 
  LastFillPrice float64 `sql:"type:DECIMAL(12,2)"`
  LastFillQuantity float64 `sql:"type:DECIMAL(12,2)"`
  RemainingQuantity float64 `sql:"type:DECIMAL(12,2)"`
  CreateDate time.Time
  TransactionDate time.Time
  Class string
  NumLegs int
  PositionReviewed string `sql:"not null;type:ENUM('No', 'Yes');default:'No'"`
  Legs []OrderLeg
}

type OrderLeg struct {
  Id uint `gorm:"primary_key"`
  UserId uint `sql:"not null;index:UserId"`
  OrderId uint `sql:"not null;index:OrderId"`  
  CreatedAt time.Time
  UpdatedAt time.Time  
  Type string
  Symbol string
  OptionSymbol string
  Side string
  Qty int
  Status string
  Duration string
  AvgFillPrice float64 `sql:"type:DECIMAL(12,2)"`
  ExecQuantity float64 `sql:"type:DECIMAL(12,2)"`
  LastFillPrice float64 `sql:"type:DECIMAL(12,2)"`
  LastFillQuantity float64 `sql:"type:DECIMAL(12,2)"`
  RemainingQuantity float64 `sql:"type:DECIMAL(12,2)"`
  CreateDate time.Time
  TransactionDate time.Time  
}
   
/* End File */