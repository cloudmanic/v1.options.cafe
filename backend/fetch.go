package main

import (
  "fmt"
	"sort"  
  "sync"
	"strings"    
  "encoding/json"
  "./brokers/types"
  "./brokers/tradier"  
)

type Fetch struct {
  mu sync.Mutex
  activeSymbols []string
  broker tradier.Api
}

//
// Return active symbols. This is handy because we 
// sort and filter before returning.
//
func (t *Fetch) GetActiveSymbols() ([]string) {
  
  // Lock da memory
	t.mu.Lock()
	defer t.mu.Unlock()   
  
  // Symbols we always want
  t.activeSymbols = append(t.activeSymbols, "$DJI")
  t.activeSymbols = append(t.activeSymbols, "SPX")
  t.activeSymbols = append(t.activeSymbols, "COMP")   
  t.activeSymbols = append(t.activeSymbols, "VIX") 
    
  // Clean up the list.
  t.activeSymbols = t.ToUpperStrings(t.activeSymbols)
  t.activeSymbols = t.RemoveDupsStrings(t.activeSymbols)
  
  // Sort the list.
  sort.Strings(t.activeSymbols)
  
  // Return the cleaned up list.
  return t.activeSymbols;
  
}

// ----------------- Market Status ------------------- //

//
// Do get watchlists
//
func (t *Fetch) GetMarketStatus() (error) {

  // Make API call
  marketStatus, err := t.broker.GetMarketStatus()
  
  if err != nil {
    return err  
  }   
  
  // Send up websocket.
  err = t.SendWebSocket("MarketStatus:refresh", marketStatus)
  
  if err != nil {
    return fmt.Errorf("GetMarketStatus() SendWebSocket : ", err)
  }   
  
  // Return Happy
  return nil
  
}

// ----------------- User Profile ------------------- //

//
// Do get user profile
//
func (t *Fetch) GetUserProfile() (error) {

  // Make API call
  userProfile, err := t.broker.GetUserProfile()
  
  if err != nil {
    return err  
  }   
  
  // Send up websocket.
  err = t.SendWebSocket("UserProfile:refresh", userProfile)
  
  if err != nil {
    return fmt.Errorf("GetUserProfile() SendWebSocket : ", err)
  }   
  
  // Return Happy
  return nil
  
}

// ----------------- Orders ------------------- //

//
// Do get user profile
//
func (t *Fetch) GetOrders() (error) {

  var orders []types.Order

  // Make API call
  orders, err := t.broker.GetOrders()
  
  if err != nil {
    return err  
  }   
  
  // Send up websocket.
  err = t.SendWebSocket("Orders:refresh", orders)
  
  if err != nil {
    return fmt.Errorf("GetUserProfile() GetOrders : ", err)
  }   
  
  // Return Happy
  return nil
  
}

// ----------------- Quotes ------------------- //

//
// Start the streaming quotes.
//
func (t *Fetch) StartStreamingQuotes(websocketSendQuoteChannel chan string) {
  
  t.broker.DoQuotes(websocketSendQuoteChannel)
  
}

//
// Do get quotes - more details from the streaming - activeSymbols
//
func (t *Fetch) GetActiveSymbolsDetailedQuotes() (error) {
   
  symbols := t.GetActiveSymbols()
  detailedQuotes, err := t.broker.GetQuotes(symbols)
  
  if err != nil {
    fmt.Println(err)
    return err
  }   
  
  // Loop through the quotes sending them up the websocket channel
  for _, row := range detailedQuotes{

    // Convert to a json string.
    data_json, err := json.Marshal(row)
    
    if err != nil {
      fmt.Println("DoDetailedQuotes() json.Marshal : ", err)
      return err
    }     

    // Send data up websocket.
    send_json, err := ws.GetSendJson("DetailedQuotes:refresh", string(data_json))  
    
    if err != nil {
      fmt.Println("DoDetailedQuotes() GetSendJson Send Object : ", err)
      return err
    } 
    
    // Send up websocket.
    err = t.SendQuoteWebSocket(send_json)
    
    if err != nil {
      return fmt.Errorf("GetActiveSymbolsDetailedQuotes() SendWebSocket : ", err)
    }     
    
  } 
  
  // Return happy
  return nil
  
}

// ----------------- Watchlists ------------------- //

//
// Do get watchlists
//
func (t *Fetch) GetWatchlists() (error) {

  watchlists, err := t.broker.GetWatchLists()
  
  if err != nil {
    return err  
  }  
  
  // Loop through and send data up websocket
  t.mu.Lock()
  
  for _, row := range watchlists {
  
    // Update active symbols 
    for _, row2 := range row.Symbols {
      t.activeSymbols = append(t.activeSymbols, row2.Name)
    }   
   
    // Send up websocket.
    err = t.SendWebSocket("Watchlist:refresh", row)
    
    if err != nil {
      return fmt.Errorf("GetWatchlists() SendWebSocket : ", err)
    }    
    
  }
  
  t.mu.Unlock() 
  
  // Return Happy
  return nil
  
}

//
// Do get a watchlist
//
func (t *Fetch) GetWatchlist(listName string) (error) {
	  
  // Get default watch list.
  watchlist, err := t.broker.GetWatchList(listName)
    
  if err != nil {
    return fmt.Errorf("GetWatchlist() : ", err)
  } 
  
  // Update our list of active symbols.
  t.activeSymbols = make([]string, 0)
  
  t.mu.Lock()
  
  for _, row := range watchlist.Symbols {
    t.activeSymbols = append(t.activeSymbols, row.Name)
  }
  
  t.mu.Unlock()
  
  // Send up websocket.
  err = t.SendWebSocket("Watchlist:refresh", watchlist)
 
  if err != nil {
    return fmt.Errorf("GetWatchlist() SendWebSocket : ", err)
  }   
  
  // Return Happy
  return nil    
  
}

// ----------------- Helper Functions ---------------- //

//
// Send data up websocket. 
//
func (t *Fetch) SendWebSocket(send_type string, sendObject interface{}) (error) {
  
  // Convert to a json string.
  dataJson, err := json.Marshal(sendObject)

  if err != nil {
    return fmt.Errorf("SendWebSocket() json.Marshal : ", err)
  } 
  
  // Send data up websocket.
  sendJson, err := ws.GetSendJson(send_type, string(dataJson))  
  
  if err != nil {
    return fmt.Errorf("SendWebSocket() GetSendJson Send Object : ", err)
  }   

  // Write data out websocket
  websocketSendChannel <- sendJson
  
  // Return happy
  return nil
  
}

//
// Send data up quote websocket. 
//
func (t *Fetch) SendQuoteWebSocket(sendJson string) (error) {

  // Write data out websocket
  websocketSendQuoteChannel <- sendJson
  
  // Return happy
  return nil
  
}

//
// Remove duplicates from an array of strings.
//
func (t *Fetch) RemoveDupsStrings(list []string) ([]string) {
	
	u := make([]string, 0, len(list))
	m := make(map[string]bool)
	
	for _, row := range list {
		
		if _, ok := m[row]; !ok {
			m[row] = true
			u = append(u, row)
		}		
		
	}

	return u
} 


//
// Take an array of strings and make them all upper case.
//
func (t *Fetch) ToUpperStrings(vs []string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = strings.ToUpper(v)
    }
    return vsm
}

/* End File */