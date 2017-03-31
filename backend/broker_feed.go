package main

import (
  "fmt"
  "time"
  "./brokers/types"
  "./brokers/tradier"  
)

type BrokerFeed struct {
  fetch Fetch   
  detailedQuotes []tradier.Quote
}

//
// When we have a broker access token, and the software lisc we call this.
// We start fetching data from the broker and such. This continues to run
// until the broker access token stops working, or their software lisc
// expires or is revoked. While the first websocket connection tiggers this.
// To go routines within this continue to run even if all clients are disconnected. 
//
func (t *BrokerFeed) Start() {
  
  // Setup tickers for broker polling.
  go t.DoOrdersTicker()
  go t.DoUserProfileTicker()
  go t.DoGetWatchlistsTicker()
  go t.DoGetDetailedQuotes()
  go t.DoGetMarketStatusTicker()
  go t.DoGetBalancesTicker()
  
  // Do Archive Calls
  go t.DoOrdersArchive()
   
}

// ---------------------- Tickers (polling) ---------------------------- //

//
// Ticker - User Profile : 60 seconds
//
func (t *BrokerFeed) DoUserProfileTicker() {

  var err error

  for {
    
    err =  t.fetch.GetUserProfile()

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
func (t *BrokerFeed) DoOrdersArchive() {

  var err error
  var orders []types.Order

  for {
    
    // Load up all orders 
    orders, err = t.fetch.GetAllOrders()
        
    if err != nil {
      fmt.Println(err)
    }       

    // Store the orders in our database
    archiveFeed.StoreOrders(orders)

    // Clear memory
    orders = nil
     
    // Sleep for 24 hours
    time.Sleep(time.Hour * 24)
        
  } 

}

//
// Ticker - Orders : 3 seconds
//
func (t *BrokerFeed) DoOrdersTicker() {

  var err error

  for {
    
    // Load up orders 
    err = t.fetch.GetOrders()
    
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
func (t *BrokerFeed) DoGetWatchlistsTicker() {

  for {
    
    // Load up our watchlists    
    err := t.fetch.GetWatchlists()
  
    if err != nil {
      fmt.Println(err)
    }
    
    // Update any active symbols
    symbols := t.fetch.GetActiveSymbols()
    t.fetch.broker.SetActiveSymbols(symbols)
    
    // Sleep for 30 second.
    time.Sleep(time.Second * 30)
        
  } 

}

//
// Ticker - Get DetailedQuotes : 1 second
//
func (t *BrokerFeed) DoGetDetailedQuotes() {

  for {
    
    // Load up our DetailedQuotes  
    err := t.fetch.GetActiveSymbolsDetailedQuotes()
  
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
func (t *BrokerFeed) DoGetMarketStatusTicker() {

  var err error

  for {
       
    // Load up market status. 
    err = t.fetch.GetMarketStatus()
    
    if err != nil {
      fmt.Println(err)
    }       
    
    // Sleep for 10 second.
    time.Sleep(time.Second * 5)
     
  }
  
}

// Ticker - Get GetBalances : 5 seconds
//
func (t *BrokerFeed) DoGetBalancesTicker() {

  var err error

  for {
       
    // Load up market status. 
    err = t.fetch.GetBalances()
    
    if err != nil {
      fmt.Println(err)
    }       
    
    // Sleep for 5 second.
    time.Sleep(time.Second * 5)
     
  }
  
}

/* End File */