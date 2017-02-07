package tradier

import (
  "fmt"
  "errors"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"   
)

type SessionStruct struct {
  Stream struct {
    Url string `json:"url"`        
    SessionId string `json:"sessionid"`
  }
}

type StreamQuote struct {
  Type string `json:"type"`
  Symbol string `json:"symbol"`
  Price string `json:"price"`
  Size string `json:"size"`        
}

type Quote struct {
  Type string
  Symbol string
  Size int
  Last float64
  Open float64 
  High float64 
  Low float64
  Close float64 
  PrevClose float64      
  Change float64
  ChangePercentage float64 `json:"change_percentage"`
  Volume int
  AverageVolume int `json:"average_volume"`
  LastVolume int `json:"last_volume"`     
  Description string       
}

//
// Get a quote.
//
func (t * Api) GetQuotes(symbols []string) ([]Quote, error) {
  
  // No symbols, no quotes.
	if len(symbols) == 0 {
		return nil, nil
	}  
  
  var quotes []Quote
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("GET", apiBaseUrl + "/markets/quotes?symbols=" + strings.Join(symbols, ","), nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))   
 
  res, err := client.Do(req)
      
  if err != nil {
    return nil, err 
  }        
  
  // Close Body
  defer res.Body.Close()  
  
  //fmt.Println(res.Header.Get("X-Ratelimit-Allowed"))  
  //fmt.Println(res.Header.Get("X-Ratelimit-Used"))
  //fmt.Println(res.Header.Get("X-Ratelimit-Available"))
  //fmt.Println(res.Header.Get("X-Ratelimit-Expiry"))

  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return nil, errors.New(fmt.Sprint("GetQuotes API did not return 200, It returned ", res.StatusCode))
  }   
    
  // Read the data we got.
  body, err := ioutil.ReadAll(res.Body)
  
  if err != nil {
    return nil, err 
  }   
  
  // Did we get one quote or many?
  if strings.Contains(string(body), "[") {
    
    var res map[string]map[string][]Quote 
    
    err := json.Unmarshal(body, &res)
		
		if err != nil {
			return nil, fmt.Errorf("%s: %s", err, body)
		}
		
		quote, ok := res["quotes"]["quote"]
		
		if ! ok {
			return nil, nil
		}		
   
    quotes = quote
    
  } else
  {
    
    var res map[string]map[string]Quote 
    
    err := json.Unmarshal(body, &res)
		
		if err != nil {
			return nil, fmt.Errorf("%s: %s", err, body)
		}
		
		quote, ok := res["quotes"]["quote"]
		
		if ! ok {
			return nil, nil
		}		
   
    quotes = []Quote{quote}
    
  }
  
  // Return happy
  return quotes, nil
}

/* End File */