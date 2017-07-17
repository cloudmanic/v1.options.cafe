package feed

import (
  "fmt"
  "sync"
  "time"
  "app.options.cafe/backend/models"  
  "app.options.cafe/backend/brokers"
  "app.options.cafe/backend/controllers"
  "app.options.cafe/backend/brokers/types"
)

type Base struct {
  User models.User
  Api brokers.Api
  
  DataChan chan controllers.SendStruct
  QuoteChan chan controllers.SendStruct
  
  muOrders sync.Mutex
  Orders []types.Order

  muWatchlists sync.Mutex
  Watchlists []types.Watchlist

  muBalances sync.Mutex
  Balances []types.Balance  

  muMarketStatus sync.Mutex
  MarketStatus types.MarketStatus

  muUserProfile sync.Mutex
  UserProfile types.UserProfile
}

type SendStruct struct {
  Type string `json:"type"`
  Data string `json:"data"`
}

//
// When we have a broker access token and an active user we call this.
// We start fetching data from the broker and such. This continues to run
// until the broker access token stops working, or the user expires
// or is revoked.
//
func (t *Base) Start() {
  
  fmt.Println("Starting Polling....")

  // Setup tickers for broker polling.
  go t.DoOrdersTicker()
  go t.DoUserProfileTicker()
  go t.DoGetWatchlistsTicker()
  go t.DoGetDetailedQuotes()
  go t.DoGetMarketStatusTicker()
  go t.DoGetBalancesTicker()
  
  // Do Archive Calls
  //go t.DoOrdersArchive()
   
}


// ---------------------- Tickers (polling) ---------------------------- //

//
// Ticker - User Profile : 60 seconds
//
func (t *Base) DoUserProfileTicker() {

  var err error

  for {
    
    err =  t.GetUserProfile()

    if err != nil {
      fmt.Println(err)
    } 
    
    // Sleep for 60 second.
    time.Sleep(time.Second * 60)
     
  }

}

//
// Ticker - Orders Archive : 24 hours
//
func (t *Base) DoOrdersArchive() {



/*
  //var positions = &[]models.Position{}
  //db.Where("user_id = ? AND trade_group_id = ?", t.userId, 132).Find(positions)  

  //archive.ClassifyTradeGroup(positions)
  


  

  var err error
  var orders []types.Order

  for {
    
    // Load up all orders 
    orders, err = t.GetAllOrders()
        
    if err != nil {
      fmt.Println(err)
    }       

    // Store the orders in our database
    archive.StoreOrders(db, orders, t.userId)

    // Clear memory
    orders = nil
     
    // Sleep for 24 hours
    time.Sleep(time.Hour * 24)
        
  } 
*/

}

//
// Ticker - Orders : 3 seconds
//
func (t *Base) DoOrdersTicker() {

  var err error

  for {
    
    // Load up orders 
    err = t.GetOrders()
    
    if err != nil {
      fmt.Println(err)
    }       
    
    // Sleep for 3 second.
    time.Sleep(time.Second * 3)
        
  } 

}

//
// Ticker - Watchlists : 30 seconds
//
func (t *Base) DoGetWatchlistsTicker() {

  for {
    
    // Load up our watchlists    
    err := t.GetWatchlists()
  
    if err != nil {
      fmt.Println(err)
    }
    
    // Update any active symbols
    symbols := t.GetActiveSymbols()
    t.Api.SetActiveSymbols(symbols)
    
    // Sleep for 30 second.
    time.Sleep(time.Second * 30)
        
  } 

}

//
// Ticker - Get DetailedQuotes : 1 second
//
func (t *Base) DoGetDetailedQuotes() {

  for {
    
    // Load up our DetailedQuotes  
    err := t.GetActiveSymbolsDetailedQuotes()
  
    if err != nil {
      fmt.Println(err)
    }
    
    // Sleep for 1 second
    time.Sleep(time.Second * 1)
        
  } 

}

//
// Ticker - Get GetMarketStatus : 5 seconds
//
func (t *Base) DoGetMarketStatusTicker() {

  var err error

  for {
       
    // Load up market status. 
    err = t.GetMarketStatus()
    
    if err != nil {
      fmt.Println(err)
    }       
    
    // Sleep for 10 second.
    time.Sleep(time.Second * 5)
     
  }
  
}

// Ticker - Get GetBalances : 5 seconds
//
func (t *Base) DoGetBalancesTicker() {

  var err error

  for {
       
    // Load up market status. 
    err = t.GetBalances()
    
    if err != nil {
      fmt.Println(err)
    }       
    
    // Sleep for 5 second.
    time.Sleep(time.Second * 5)
     
  }
  
}

/* End File */