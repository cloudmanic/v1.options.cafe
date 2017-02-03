package tradier

import (
  "fmt"
  "errors"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "../types"  
)

//
// Get all watchlists
//
func (t * Api) GetWatchLists() ([]types.Watchlist, error) {
  
  var watchlists []types.Watchlist
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("GET", apiBaseUrl + "/watchlists", nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))   
 
  res, err := client.Do(req)
      
  if err != nil {
    return watchlists, err  
  }        
  
  // Close Body
  defer res.Body.Close()   
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return watchlists, errors.New(fmt.Sprint("GetWatchLists API did not return 200, It returned ", res.StatusCode)) 
  }    
     
  // Read the data we got.
  body, err := ioutil.ReadAll(res.Body)

  // Do we have one watchlist or many?
  var one map[string]map[string]*types.Watchlist
  
  if err := json.Unmarshal(body, &one); err != nil {
         
    if err := json.Unmarshal(body, &watchlists); err != nil {
      return watchlists, err       
    }
          
  } else
  {
    // Build the struct since we only got one in our return.
    watchlists = make([]types.Watchlist, 1)
    watchlists[0].Id = one["watchlists"]["watchlist"].Id
    watchlists[0].Name = one["watchlists"]["watchlist"].Name     
  }

	// If we made it this far lets make calls to fill in the watch lists.
	for key, row := range watchlists {

    // Get the complete watchlist.      
    ws, err := t.GetWatchList(row.Id)
    
    if err != nil {
      fmt.Println(err)
      continue;
    }
    	  	
  	watchlists[key].Symbols = ws.Symbols
  }
	  
  // Return data we just got. 
  return watchlists, nil
}

//
// Get a specific watchlist
//
func (t * Api) GetWatchList(id string) (types.Watchlist, error) {
  
  var watchlist types.Watchlist
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("GET", apiBaseUrl + "/watchlists/" + id, nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))   
 
  res, err := client.Do(req)
      
  if err != nil {
    return watchlist, err  
  }        
  
  // Close Body
  defer res.Body.Close()    
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return watchlist, errors.New(fmt.Sprint("GetWatchLists API did not return 200, It returned ", res.StatusCode)) 
  }    
     
  // Read the data we got.
  body, _ := ioutil.ReadAll(res.Body) 
  
  // Temp struct to hold response.
  type WatchlistResponse struct {
    Name string 
    Id string
    ItemsJson json.RawMessage `json:"items"`     
  }  
  
  // Bust open the watchlist.
  var ws map[string]*WatchlistResponse 
  
  if err := json.Unmarshal(body, &ws); err != nil {
    return watchlist, err 
  }
  
  // Setup the watchlist
  watchlist.Id = ws["watchlist"].Id
  watchlist.Name = ws["watchlist"].Name  

  // Get the items
  var items map[string][]*types.WatchlistSymbol
  
  // Do we have one item or many items
  if err := json.Unmarshal(ws["watchlist"].ItemsJson, &items); err != nil {
    
    var items map[string]*types.WatchlistSymbol
    
    if err := json.Unmarshal(ws["watchlist"].ItemsJson, &items); err != nil {
      return watchlist, err 
    }
    
    // Set a single symbol
    watchlist.Symbols = []types.WatchlistSymbol{ { Id: items["item"].Id, Name: items["item"].Name } }    
  } else
  {  
    // Loop through and add the symbols
  	for _, row := range items["item"] {
    	watchlist.Symbols = append(watchlist.Symbols, *row)    	
    }
  }  
  
  // Return happy
  return watchlist, nil	
  
}

/* End File */