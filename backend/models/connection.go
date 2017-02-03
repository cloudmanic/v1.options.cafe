package models

import (
  //"os"
  //"log"
  "github.com/jinzhu/gorm"
  //_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
  Conn *gorm.DB
} 

/*
//
// Connect to the db and run migrations.
//
func (m Connection) Connect() (*gorm.DB) {
  
  var err error
    
  // Connect to Mysql
  m.Conn, err := gorm.Open("mysql", os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@/" + os.Getenv("DB_DATABASE") + "?charset=utf8&parseTime=True&loc=Local")
  
  if err != nil {
    panic("failed to connect database")
  }

  // Enable
  m.Conn.LogMode(true)
  m.Conn.SetLogger(log.New(os.Stdout, "\r\n", 0))

  // Migrate the schemas (one per table).
  m.Conn.AutoMigrate(&User{})
  
  return m.Conn   
}
*/

/* End File */