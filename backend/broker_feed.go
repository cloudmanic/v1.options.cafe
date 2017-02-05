package main

import (
  "fmt"
  "time"
  "./brokers/tradier" 
)

type BrokerFeed struct {
  fetch Fetch   
  activeSymbols []string
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
   
}

// ---------------------- Tickers (polling) ---------------------------- //

//
// Ticker - User Profile : 5 seconds
//
func (t *BrokerFeed) DoUserProfileTicker() {

  var err error

  for {
    
    err =  t.fetch.GetUserProfile()

    if err != nil {
      fmt.Println(err)
    } 
    
    // Sleep for 5 second.
    time.Sleep(time.Second * 5)
     
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
// Ticker - Watchlists : 10 seconds
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
    
    // Sleep for 10 second.
    time.Sleep(time.Second * 10)
        
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
// Ticker - Get GetMarketStatus : 10 seconds
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
    time.Sleep(time.Second * 10)
     
  }
  
}

/* End File */