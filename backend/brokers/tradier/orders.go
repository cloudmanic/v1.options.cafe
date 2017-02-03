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
  
  orders := make([]types.Order, 0)
  var t_orders []types.TradierOrder
  
  // Setup http client
  client := &http.Client{}    
  
  // Setup api request
  req, _ := http.NewRequest("GET", apiBaseUrl + "/accounts/" + t.defaultAccountId +  "/orders", nil)
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

  // Make sure we have at least one order.
  vo := gjson.Get(string(body), "orders.order")
  
  if ! vo.Exists() {
    // Return happy (just no orders found)
    return orders, nil	
  }
  
  // See if we have no orders, one order, or many orders
  value := gjson.Get(string(body), "orders.order.id")
  
  // More than one order
  if ! value.Exists() {
    
    var ws map[string]map[string][]types.TradierOrder
    
    // Unmarshal json
    if err := json.Unmarshal(body, &ws); err != nil {
      return orders, err 
    }
    
    // Set the orders we return.
    t_orders = ws["orders"]["order"]
         
  } else 
  {    
    var ws map[string]map[string]types.TradierOrder
    
    // Unmarshal json
    if err := json.Unmarshal(body, &ws); err != nil {
      return orders, err 
    }
    
    // Set the orders we return.
    t_orders = append(t_orders, ws["orders"]["order"])   
    
  }  
  
  // Clean up the array and make it more generic
  for _, row := range t_orders {
    
    orders = append(orders, types.Order{
      Id: row.Id,
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
      NumLegs: row.NumLegs})
    
  }
  
  // Return happy
  return orders, nil	
  
}

/* End File */