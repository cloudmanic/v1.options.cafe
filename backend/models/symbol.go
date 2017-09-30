//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
  "time"
  "app.options.cafe/backend/library/services"
)

type Symbol struct {
  Id uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time 
  ShortName string `sql:"not null"`  
  Name string `sql:"not null"`
}  

//
// Create a new Symbol entry.
//
func (t * DB) CreateNewSymbol(short string, name string) (Symbol, error) {
    
  var symb Symbol
    
  // First make sure we don't already have this symbol
  if t.Connection.Where("short_name = ?", short).First(&symb).RecordNotFound() {

    // Create entry.
    symb = Symbol{ Name: name, ShortName: short }
              
    t.Connection.Create(&symb)
    
    // Log Symbol creation.
    services.Log("CreateNewSymbol - Created a new Symbol entry - (" + short + ") " + name)      

  }  
  
  // Return the user.
  return symb, nil  
   
} 
    
          
/* End File */