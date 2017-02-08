package tradier

import (
  "fmt"
  "errors"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "../types"
  "github.com/tidwall/gjson"  
)

//
// Get Orders
//
func (t * Api) GetOrders() ([]types.Order, error) {
  
  var orders []types.Order
  var t_orders []types.TradierOrder
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("GET", apiBaseUrl + "/user/orders", nil)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", fmt.Sprint("Bearer ", t.ApiKey))   
 
  res, err := client.Do(req)
      
  if err != nil {
    return orders, err  
  }        
  
  // Close Body
  defer res.Body.Close()    
  
  // Make sure the api responded with a 200
  if res.StatusCode != 200 {
    return orders, errors.New(fmt.Sprint("GetOrders API did not return 200, It returned ", res.StatusCode)) 
  }    
     
  // Read the data we got.
  body, _ := ioutil.ReadAll(res.Body) 

  // Make sure we have at least one account (this should never happen)
  vo := gjson.Get(string(body), "accounts.account")
  
  if ! vo.Exists() {
    // Return happy (just no orders found)
    return orders, nil	
  }

  // Do we have only one account?
  vo = gjson.Get(string(body), "accounts.account.account_number")
  
  // Only one account
  if vo.Exists() {

    if t.orderParseOneAccount(body, &t_orders) != nil {
      return orders, nil
    }
           
  } else // More than one accounts
  {

    if t.orderParseMoreThanOneAccount(body, &t_orders) != nil {
      return orders, nil
    }   
        
  }
  
  // Clean up the array and make it more generic
  for _, row := range t_orders {
    
    var legs []types.OrderLeg
    
    // Merge in the legs
    if len(row.Legs) > 0 {
      
      for _, row2 := range row.Legs {
        
        legs = append(legs, types.OrderLeg{
          Type: row2.Type,
          Symbol: row2.Symbol,
          OptionSymbol: row2.OptionSymbol, 
          Side: row2.Side, 
          Quantity: row2.Quantity, 
          Status: row2.Status, 
          Duration: row2.Duration, 
          AvgFillPrice: row2.AvgFillPrice, 
          ExecQuantity: row2.ExecQuantity, 
          LastFillPrice: row2.LastFillPrice, 
          LastFillQuantity: row2.LastFillQuantity, 
          RemainingQuantity: row2.RemainingQuantity, 
          CreateDate: row2.CreateDate, 
          TransactionDate: row2.TransactionDate,           
        })
        
      }
      
    }
    
    // Append the orders
    orders = append(orders, types.Order{
      Id: row.Id,
      AccountId: row.AccountId,
      Type: row.Type,
      Symbol: row.Symbol,
      Side: row.Side,
      Quantity: row.Quantity,
      Status: row.Status,
      Duration: row.Duration,
      Price: row.Price,
      AvgFillPrice: row.AvgFillPrice, 
      ExecQuantity: row.ExecQuantity, 
      LastFillPrice: row.LastFillPrice,
      LastFillQuantity: row.LastFillQuantity,
      RemainingQuantity: row.RemainingQuantity,
      CreateDate: row.CreateDate,
      TransactionDate: row.TransactionDate,
      Class: row.Class,
      NumLegs: row.NumLegs,
      Legs: legs,
    })
    
  }
  
  // Return happy
  return orders, nil	
  
}

//
// Parse the case where the user just has more than one account.
// This function just sets the t_orders
//
func (t * Api) orderParseMoreThanOneAccount(body []byte, t_orders *[]types.TradierOrder) error {
    
  vo := gjson.Get(string(body), "accounts.account")
  
  // Loop through the different accounts.
  vo.ForEach(func(key, value gjson.Result) bool{
    
    // Do we have any orders.
    vo2 := gjson.Get(value.String(), "orders")
    
    if ! vo2.Exists() {
      return true
    }
    
    // Set the account id
    account_number := gjson.Get(value.String(), "account_number").String() 
    
    // Do we have more than one order
    vo2 = gjson.Get(value.String(), "orders.order.id")
    
    // More than one order??
    if ! vo2.Exists() {
              
      var ws []types.TradierOrder
      
      // Get just the orders part
      vo3 := gjson.Get(value.String(), "orders.order")
      
      // Unmarshal json
      if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
        return true 
      }
      
      // Set the orders to our return
      for _, row := range ws {        
        row.AccountId = account_number
        *t_orders = append(*t_orders, row)
      }
           
    } else 
    {    
      var ws types.TradierOrder
  
      // Get just the orders part
      vo3 := gjson.Get(value.String(), "orders.order")
  
      // Unmarshal json
      if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
        return true 
      }

      // Set the orders we return.
      ws.AccountId = account_number
      *t_orders = append(*t_orders, ws)   
      
    }       
    
    // keep iterating
    return true 
    
  }) 

  // Happy
  return nil      
}

//
// Parse the case where the user just has one account.
// This function just sets the t_orders
//
func (t * Api) orderParseOneAccount(body []byte, t_orders *[]types.TradierOrder) error {
  
  // Do we have any orders.
  vo2 := gjson.Get(string(body), "accounts.account.orders")
  
  if ! vo2.Exists() {
    return errors.New("No Orders Found") 
  }
  
  // Set the account id
  account_number := gjson.Get(string(body), "accounts.account.account_number").String() 
  
  // Do we have more than one order
  vo2 = gjson.Get(string(body), "accounts.account.orders.order.id")
  
  // More than one order??
  if ! vo2.Exists() {
            
    var ws []types.TradierOrder
    
    // Get just the orders part
    vo3 := gjson.Get(string(body), "accounts.account.orders.order")
    
    // Unmarshal json
    if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
      return err
    }
    
    // Set the orders to our return
    for _, row := range ws {
      row.AccountId = account_number
      *t_orders = append(*t_orders, row)
    }
  
  } else 
  {    
    var ws types.TradierOrder
  
    // Get just the orders part
    vo3 := gjson.Get(string(body), "accounts.account.orders.order")
  
    // Unmarshal json
    if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
      return err
    }
          
    // Set the orders we return.
    ws.AccountId = account_number
    *t_orders = append(*t_orders, ws)   
    
  }      
  
  // Success
  return nil
}

/* End File */