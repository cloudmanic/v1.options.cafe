package main

import (
  "fmt"
  "./brokers/types"   
)

type ArchiveFeed struct {}

//
// Pass in all orders and archive them by putting them into our database.
//
func (t * ArchiveFeed) StoreOrders(orders []types.Order) (error) {
  
  // Loop through the orders and process
  for _, row := range orders {
      
    fmt.Println("Account: " + row.AccountId)
    
  }
  
  // Return Happy
  return nil
  
}

/* End File */