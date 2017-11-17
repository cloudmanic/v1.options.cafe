package tradier

import (
  "crypto/md5"
  "encoding/hex"
  "github.com/tidwall/gjson"
  "github.com/app.options.cafe/backend/brokers/types"    
)

//
// Get History by Account Id
//
func (t * Api) GetHistoryByAccountId(accountId string) ([]types.History, error) {
  
  var history []types.History
  
  // Get the JSON - Set big limit to get all data
  jsonRt, err := t.SendGetRequest("/accounts/" + accountId + "/history?limit=1000000")

  if err != nil {
    return history, err  
  } 

  // Make sure we have at least one event
  vo := gjson.Get(jsonRt, "history.event")
  
  if ! vo.Exists() {
    return history, nil	
  }  
 
  // Do we have only one event or more?
  vo = gjson.Get(jsonRt, "history.event.date")
 
  // Only one event
  if vo.Exists() {
      
    // Add to balances array
    history = append(history, t.addJsonToEvent(gjson.Get(jsonRt, "history.event").String()))
           
  } else // More than one historical event
  {

    vo := gjson.Get(jsonRt, "history.event")
  
    // Loop through the different accounts.
    vo.ForEach(func(key, value gjson.Result) bool {
      
      // Add to balances array
      history = append(history, t.addJsonToEvent(value.String()))
      
      // keep iterating
      return true  
          
    }) 
        
  } 
  
  // Return happy
  return history, nil
  
}

//
// Take a json string and turn it into a history object.
//
func (t * Api) addJsonToEvent(eventJson string) types.History {
    
  hasher := md5.New()
  hasher.Write([]byte(eventJson))
    
  // Return History event.
  return types.History{
    BrokerId: hex.EncodeToString(hasher.Sum(nil)),
    Type: gjson.Get(eventJson, "type").String(),
    Date: gjson.Get(eventJson, "date").String(),
    Amount: gjson.Get(eventJson, "amount").Float(),
    Symbol: gjson.Get(eventJson, "trade.symbol").String(),
    Commission: gjson.Get(eventJson, "trade.commission").Float(),
    Description: gjson.Get(eventJson, "trade.description").String(),
    Price: gjson.Get(eventJson, "trade.price").Float(),
    Quantity: gjson.Get(eventJson, "trade.quantity").Int(), 
    TradeType: gjson.Get(eventJson, "trade.trade_type").String(),         
  }
  
}

/* End File */