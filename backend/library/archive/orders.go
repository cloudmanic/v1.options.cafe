package archive

import (
  "fmt"
  "time"
  "github.com/jinzhu/gorm"  
  "github.com/stvp/rollbar"
  "github.com/app.options.cafe/backend/models"
  "github.com/app.options.cafe/backend/brokers/types"     
)

//
// Pass in all orders and archive them by putting them into our database.
// We only archive orders that are filled. TODO: Review partially_filled orders.
//
func StoreOrders(db *gorm.DB, orders []types.Order, userId uint) (error) {
  
  // Loop through the orders and process
  for _, row := range orders {
      
    // We only care about filled orders
    if row.Status != "filled" {
      continue;
    }
    
    // See if we already have this record in our database
    var count int
    order := &models.Order{}
    
    db.Where("broker_id = ? AND user_id = ?", row.Id, userId).First(order).Count(&count)
    
    if count > 0 {
      continue
    }
      
    // Timestamp Layout
    layout := "2006-01-02T15:04:05.000Z"
      
    // Convert Create Date
    createDate, err := time.Parse(layout, row.CreateDate)

    if err != nil {
      fmt.Println(err)
      rollbar.Error(rollbar.ERR, err)
      continue
    }

    // Convert TransactionDate
    transactionDate, err := time.Parse(layout, row.TransactionDate)

    if err != nil {
      fmt.Println(err)
      rollbar.Error(rollbar.ERR, err)
      continue
    }
      
    // Create object we insert into the DB
    order = &models.Order{
                UserId: userId,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
                BrokerId: row.Id,
                AccountId: row.AccountId,  
                Type: row.Type,
                Symbol: row.Symbol,
                Side: row.Side,
                Qty: int(row.Quantity),
                Status: row.Status,
                Duration: row.Duration,
                Price: row.Price,
                AvgFillPrice: row.AvgFillPrice,  
                ExecQuantity: row.ExecQuantity,  
                LastFillPrice: row.LastFillPrice, 
                LastFillQuantity: row.LastFillQuantity,
                RemainingQuantity: row.RemainingQuantity,
                CreateDate: createDate, 
                TransactionDate: transactionDate,
                Class: row.Class,
                PositionReviewed: "No",
                NumLegs: row.NumLegs,        
              }
    
    // Insert into DB          
    db.Create(&order)
     
    // Loop through the order legs and add them
    for _, row2 := range row.Legs {
      
      // Convert Create Date
      createDate, err := time.Parse(layout, row2.CreateDate)
      
      if err != nil {
        fmt.Println(err)
        rollbar.Error(rollbar.ERR, err)
        continue
      }
      
      // Convert TransactionDate
      transactionDate, err := time.Parse(layout, row2.TransactionDate)
      
      if err != nil {
        fmt.Println(err)
        rollbar.Error(rollbar.ERR, err)
        continue
      }
            
      // Create object we insert into the DB
      leg := &models.OrderLeg{
                  UserId: userId,
                  OrderId: order.Id,  
                  CreatedAt: time.Now(),
                  UpdatedAt: time.Now(),  
                  Type: row2.Type,
                  Symbol: row2.Symbol,
                  OptionSymbol: row2.OptionSymbol,
                  Side: row2.Side,
                  Qty: int(row2.Quantity),
                  Status: row2.Status,
                  Duration: row2.Duration,
                  AvgFillPrice: row2.AvgFillPrice,
                  ExecQuantity: row2.ExecQuantity,
                  LastFillPrice: row2.LastFillPrice,
                  LastFillQuantity: row2.LastFillQuantity,
                  RemainingQuantity: row2.RemainingQuantity,
                  CreateDate: createDate,
                  TransactionDate: transactionDate,        
              }
      
      // Insert into DB          
      db.Create(&leg)
      
    }
    
  }
  
  // Now build out our positions database table based on past orders.
  StorePositions(db, userId)
   
  // Return Happy
  return nil
  
}

/* End File */