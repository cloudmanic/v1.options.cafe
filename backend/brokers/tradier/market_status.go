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
// Get market status
//
func (t * Api) GetMarketStatus() (types.MarketStatus, error) {
  
  var status types.MarketStatus
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("GET", apiBaseUrl + "/markets/clock", nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))   
 
  res, err := client.Do(req)
      
  if err != nil {
    return status, err  
  }        
  
  // Close Body
  defer res.Body.Close()    
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return status, errors.New(fmt.Sprint("GetMarketStatus API did not return 200, It returned ", res.StatusCode)) 
  }    
     
  // Read the data we got.
  body, _ := ioutil.ReadAll(res.Body) 
  
  // Bust open the watchlist.
  var ws map[string]types.MarketStatus 
  
  if err := json.Unmarshal(body, &ws); err != nil {
    return status, err 
  }
  
  // Set the status we return.
  status = ws["clock"]
  
  // Return happy
  return status, nil	
  
}

/* End File */