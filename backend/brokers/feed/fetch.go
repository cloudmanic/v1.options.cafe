package feed

import (
  "fmt"
	"sort"
	"strings"    
  "encoding/json"
  "app.options.cafe/backend/controllers"  
  "app.options.cafe/backend/brokers/types"
  //"app.options.cafe/backend/library/services"   
)

//
// Send cached data up the websocket pipeline. 
// This is useful for when the application wants
// data between pooling of the broker.
//
func (t *Base) RefreshFromCached() error {
  
  // UserProfile - Send up websocket.
  err := t.WriteDataChannel("UserProfile:refresh", t.UserProfile)
  
  if err != nil {
    return fmt.Errorf("RefreshFromCached() WriteDataChannel - UserProfile:refresh : ", err)
  }   
  
  // MarketStatus - Send up websocket.
  err = t.WriteDataChannel("MarketStatus:refresh", t.MarketStatus)
  
  if err != nil {
    return fmt.Errorf("RefreshFromCached() WriteDataChannel - MarketStatus:refresh : ", err)
  } 
  
  // Orders - Send up websocket.
  err = t.WriteDataChannel("Orders:refresh", t.Orders)
  
  if err != nil {
    return fmt.Errorf("RefreshFromCached() WriteDataChannel - Orders:refresh : ", err)
  }
  
  // Balances - Send up websocket.
  err = t.WriteDataChannel("Balances:refresh", t.Balances)
  
  if err != nil {
    return fmt.Errorf("RefreshFromCached() WriteDataChannel - Balances:refresh : ", err)
  }     
  
  // Watchlists - Loop through and send data up websocket  
  for _, row := range t.Watchlists {
     
    // Send up websocket.
    err = t.WriteDataChannel("Watchlist:refresh", row)
    
    if err != nil {
      return fmt.Errorf("RefreshFromCached() WriteDataChannel - Watchlist:refresh : ", err)
    }    
    
  }  
  
  // No error
  return nil
}

//
// Return active symbols. This is handy because we 
// sort and filter before returning.
//
func (t *Base) GetActiveSymbols() []string {
  
  var activeSymbols []string  
  
  // Symbols we always want
  activeSymbols = append(activeSymbols, "$DJI")
  activeSymbols = append(activeSymbols, "SPX")
  activeSymbols = append(activeSymbols, "COMP")   
  activeSymbols = append(activeSymbols, "VIX")
  
  // Watchlists to the active symbols.
  t.muWatchlists.Lock()
  
  for _, row := range t.Watchlists {
    
    for _, row2 := range row.Symbols {
    
      activeSymbols = append(activeSymbols, row2.Name)
  
    }
    
  }
  
  t.muWatchlists.Unlock()    
  
  // Add in the orders we want.
  t.muOrders.Lock()
  
  for _, row := range t.Orders {
    
    activeSymbols = append(activeSymbols, row.Symbol)
    
    if row.NumLegs > 0 {
      
      for _, row2 := range row.Legs {
        
        activeSymbols = append(activeSymbols, row2.OptionSymbol);
      
      } 
   
    }
  
  }
  
  t.muOrders.Unlock() 
    
  // Clean up the list.
  activeSymbols = t.ToUpperStrings(activeSymbols)
  activeSymbols = t.RemoveDupsStrings(activeSymbols)
  
  // Sort the list.
  sort.Strings(activeSymbols)
  
  // Return the cleaned up list.
  return activeSymbols;
  
}

// ----------------- Market Status ------------------- //

//
// Do get watchlists
//
func (t *Base) GetMarketStatus() (error) {

  // Make API call
  marketStatus, err := t.Api.GetMarketStatus()
  
  if err != nil {
    return err  
  }   
  
  // Save the market status in the fetch object
  t.muMarketStatus.Lock()
  t.MarketStatus = marketStatus
  t.muMarketStatus.Unlock()   
  
  // Send up websocket.
  err = t.WriteDataChannel("MarketStatus:refresh", marketStatus)
  
  if err != nil {
    return fmt.Errorf("GetMarketStatus() WriteDataChannel : ", err)
  }   
  
  // Return Happy
  return nil
  
}

// ----------------- User Profile ------------------- //

//
// Do get user profile
//
func (t *Base) GetUserProfile() (error) {

  // Make API call
  userProfile, err := t.Api.GetUserProfile()
  
  if err != nil {
    return err  
  }   
  
  // Save the orders in the fetch object
  t.muUserProfile.Lock()
  t.UserProfile = userProfile
  t.muUserProfile.Unlock()    
  
  // Send up websocket.
  err = t.WriteDataChannel("UserProfile:refresh", userProfile)
  
  if err != nil {
    return fmt.Errorf("GetUserProfile() WriteDataChannel : ", err)
  }   
  
  // Return Happy
  return nil
  
}

// ----------------- Orders ------------------- //

//
// Do get orders
//
func (t *Base) GetOrders() (error) {

  var orders []types.Order

  // Make API call
  orders, err := t.Api.GetOrders()
 
  if err != nil {
    return err  
  }    
  
  // Save the orders in the fetch object
  t.muOrders.Lock()
  t.Orders = orders
  t.muOrders.Unlock() 
  
  // Send up websocket.
  err = t.WriteDataChannel("Orders:refresh", orders)
  
  if err != nil {
    return fmt.Errorf("Fetch.GetOrders() : ", err)
  }   
  
  // Return Happy
  return nil
  
}

//
// Do get all orders. We return the orders instead of sending it up the websocket
//
func (t *Base) GetAllOrders() ([]types.Order, error) {

  var orders []types.Order

  // Make API call
  orders, err := t.Api.GetAllOrders()

  if err != nil {
    return orders, fmt.Errorf("Fetch.GetAllOrders() : ", err)
  }   
  
  // Return Happy
  return orders, nil
  
}

// ----------------- Quotes ------------------- //


//
// Do get quotes - more details from the streaming - activeSymbols
//
func (t *Base) GetActiveSymbolsDetailedQuotes() (error) {
  
  symbols := t.GetActiveSymbols()  
  detailedQuotes, err := t.Api.GetQuotes(symbols)
  
  if err != nil {
    fmt.Println("DoDetailedQuotes() t.Api.GetQuotes : ", err)
    return err
  }   
  
  // Loop through the quotes sending them up the websocket channel
  for _, row := range detailedQuotes {

    // Convert to a json string.
    data_json, err := json.Marshal(row)
    
    if err != nil {
      fmt.Println("DoDetailedQuotes() json.Marshal : ", err)
      return err
    }     

    // Send data up websocket.
    send_json, err := t.GetSendJson("DetailedQuotes:refresh", string(data_json))  
    
    if err != nil {
      fmt.Println("DoDetailedQuotes() GetSendJson Send Object : ", err)
      return err
    } 
    
    // Send up websocket.
    err = t.WriteQuoteChannel(send_json)
    
    if err != nil {
      return fmt.Errorf("GetActiveSymbolsDetailedQuotes() WriteDataChannel : ", err)
    }     
    
  } 
  
  // Return happy
  return nil
  
}

// ----------------- Watchlists ------------------- //

//
// Do get watchlists
//
func (t *Base) GetWatchlists() (error) {

  watchlists, err := t.Api.GetWatchLists()
  
  if err != nil {
    return err  
  } 
  
  // Save the watchlists in the fetch object
  t.muWatchlists.Lock()
  t.Watchlists = watchlists
  t.muWatchlists.Unlock()   
  
  // Loop through and send data up websocket  
  for _, row := range watchlists {
     
    // Send up websocket.
    err = t.WriteDataChannel("Watchlist:refresh", row)
    
    if err != nil {
      return fmt.Errorf("GetWatchlists() WriteDataChannel : ", err)
    }    
    
  }
  
  // Return Happy
  return nil
  
}

// ----------------- Balances ------------------- //

//
// Do get Balances
//
func (t *Base) GetBalances() (error) {

  balances, err := t.Api.GetBalances()
  
  if err != nil {
    return err  
  } 
  
  // Save the balances in the fetch object
  t.muBalances.Lock()
  t.Balances = balances
  t.muBalances.Unlock()   
  
  // Send up websocket.
  err = t.WriteDataChannel("Balances:refresh", balances)
  
  if err != nil {
    return fmt.Errorf("GetBalances() WriteDataChannel : ", err)
  } 
  
  // Return Happy
  return nil
  
}

// ----------------- Access Tokens ------------------- //

//
// Do update access token from refresh
//
func (t *Base) AccessTokenRefresh() (error) {

  err := t.Api.DoRefreshAccessTokenIfNeeded(t.User)
  
  if err != nil {
    return err  
  } 

  // Return Happy
  return nil
  
}

// ----------------- Helper Functions ---------------- //

//
// Return a json object ready to be sent up the websocket
//
func (t *Base) GetSendJson(send_type string, data_json string) (string, error) {
  
  // Send Object
  send := SendStruct{
    Type: send_type,
    Data: data_json,
  } 
  
  send_json, err := json.Marshal(send)

  if err != nil {
    return "", err
  } 

  return string(send_json), nil
}

//
// Send data up websocket. 
//
func (t *Base) WriteDataChannel(send_type string, sendObject interface{}) (error) {

  // Convert to a json string.
  dataJson, err := json.Marshal(sendObject)

  if err != nil {
    return fmt.Errorf("WriteDataChannel() json.Marshal : ", err)
  } 
  
  // Send data up websocket.
  sendJson, err := t.GetSendJson(send_type, string(dataJson))  
  
  if err != nil {
    return fmt.Errorf("WriteDataChannel() GetSendJson Send Object : ", err)
  }   

  // Write data out websocket
  t.DataChan <- controllers.SendStruct{ UserId: t.User.Id, Message: sendJson }
  
  // Return happy.
  return nil
}

//
// Send data up quote websocket. 
//
func (t *Base) WriteQuoteChannel(sendJson string) (error) {
    
  // Write data out websocket
  t.QuoteChan <- controllers.SendStruct{ UserId: t.User.Id, Message: sendJson }

  // Return happy.
  return nil
}

//
// Remove duplicates from an array of strings.
//
func (t *Base) RemoveDupsStrings(list []string) ([]string) {
	
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
func (t *Base) ToUpperStrings(vs []string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = strings.ToUpper(v)
    }
    return vsm
}

/* End File */