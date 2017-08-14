package models

import (
  "os"
  "log"
  "github.com/jinzhu/gorm"
  _ "github.com/go-sql-driver/mysql"
  "app.options.cafe/backend/library/services"
)

type DB struct {
  Connection *gorm.DB 
}

//
// Start the DB connection.
//
func (t * DB) Start() {
  
  var err error
    
  // Connect to Mysql
  conn, err := gorm.Open("mysql", os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + "/" + os.Getenv("DB_DATABASE") + "?charset=utf8&parseTime=True&loc=Local")
  
  if err != nil {
    services.Error(err, "Failed to connect database")
    log.Fatal(err)
  }

  // Set the connection for the struct
  t.Connection = conn

  // Enable
  //t.Connection.LogMode(true)
  //t.Connection.SetLogger(log.New(os.Stdout, "\r\n", 0))

  // Migrate the schemas (one per table).
  t.Connection.AutoMigrate(&User{})
  t.Connection.AutoMigrate(&Broker{})
  t.Connection.AutoMigrate(&Order{})
  t.Connection.AutoMigrate(&OrderLeg{})
  t.Connection.AutoMigrate(&Position{}) 
  t.Connection.AutoMigrate(&TradeGroup{}) 
  
}

/* End File */