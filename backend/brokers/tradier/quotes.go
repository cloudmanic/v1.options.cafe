package tradier

import (
  "fmt"
  "bufio"
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
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return nil, errors.New(fmt.Sprint("GetQuotes API did not return 200, It returned ", res.StatusCode))
  }   
  
  // Close Body
  defer res.Body.Close()  
  
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

//
// Get quotes from Tradier
//
func (t * Api) DoQuotes(websocket_channel chan string) {
  
  // Do nothing until we have an api key 
  t.activeSymbols = "aapl,bac,frsh,ge,hd,spx,spy,vix,vxx,xlf,$dji,comp"
  
  for {
    
    println("DoQuotes: Start")
    err := t.StartQuoteStream(websocket_channel);
    
    if err != nil {
      println(fmt.Sprint("DoQuotes: StartQuoteStream - ", err))   
    }
    
  }
     
}

//
// Start quote stream
//
func (t * Api) StartQuoteStream(websocket_channel chan string) (error) {
  
  // Get the session token for streaming.
  session_id, session_url, err := t.GetSessionIdForStreaming();

  if err != nil {
    return err  
  } 
  
  // Cache the strings this is a flag for when we break the session and start a new one.
  t.mu.Lock()
  t.oldActiveSymbols = t.activeSymbols
  t.mu.Unlock()
  
  // Start the data stream.
  client := &http.Client{}    
  
  // Setup api request
  url := fmt.Sprint(session_url, "?sessionid=", session_id, "&symbols=", t.activeSymbols, "&filter=trade&linebreak=true")
  
  req, _ := http.NewRequest("POST", url, nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))  
  
  res, err := client.Do(req)
  
  if err != nil {
    return err 
  }  
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return errors.New(fmt.Sprint("Streaming API did not return 200, It returned ", res.StatusCode))
  }      
  
  reader := bufio.NewReader(res.Body)
  
  // Just keep going until we stop getting data.
  for {
    line, err := reader.ReadBytes('\n')
    
    if err != nil {
      return err
    }     

    // Json to Object TODO: right now we do not do much with this but some day we might.
    var quote StreamQuote
            
    if err := json.Unmarshal(line, &quote); err != nil {
      return err       
    }
    
    // Send to the channel for processing.  
    websocket_channel <- string(line) 
    
    // See if we got a new symbol we need to deal with
    if t.oldActiveSymbols != t.activeSymbols {
      return errors.New("Active Symbols Have Changed.......Restarting Tradier Stream")
    }  
  }

  return nil
}

//
// Get a session id for streaming.
//
func (t * Api) GetSessionIdForStreaming() (string, string, error) {
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("POST", apiBaseUrl + "/markets/events/session", nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))   
 
  res, err := client.Do(req)
      
  if err != nil {
    return "", "", err    
  }        
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return "", "", err
  }    
     
  // Read the data we got.
  body, err := ioutil.ReadAll(res.Body) 
  
  // Json to Object
  var session SessionStruct
          
  if err := json.Unmarshal(body, &session); err != nil {
    return "", "", err        
  }    
     
  // Return data we just got. 
  return session.Stream.SessionId, session.Stream.Url, nil
}