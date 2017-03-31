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
  Quantity float64
  Status string
  Duration string
  Price float64
  AvgFillPrice float64 
  ExecQuantity float64 
  LastFillPrice float64
  LastFillQuantity float64
  RemainingQuantity float64
  CreateDate time.Time
  TransactionDate time.Time
  Class string
  NumLegs int
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
  Quantity float64
  Status string
  Duration string
  AvgFillPrice float64
  ExecQuantity float64
  LastFillPrice float64
  LastFillQuantity float64
  RemainingQuantity float64
  CreateDate time.Time
  TransactionDate time.Time  
}
   
/* End File */