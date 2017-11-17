package tradier

import (
  "errors"
  "encoding/json"
  "github.com/tidwall/gjson"
  "github.com/app.options.cafe/backend/brokers/types"    
)

//
// Get All Orders Ever
//
func (t * Api) GetAllOrders() ([]types.Order, error) {
  
  var orders []types.Order

  // Get the user profile
  user, err := t.GetUserProfile()

  if err != nil {
    return orders, err  
  } 

  // Loop through the different acconts and get all orders.
  for _, row := range user.Accounts {
    
    err := t.processOrderDataForGetAllOrders(row.AccountNumber, &orders)
    
    if err != nil {
      return orders, err  
    }

  }
  
  // Return the data.
  return orders, nil  
  
}

//
// Get Orders
//
func (t * Api) GetOrders() ([]types.Order, error) {
  
  var orders []types.Order
  
  // Get the JSON
  jsonRt, err := t.SendGetRequest("/user/orders")

  if err != nil {
    return orders, err  
  } 
  
  // Process JSON returned
  return t.parseOrdersJson(jsonRt) 
  
}

//
// Process JSON returned from an ORDERS api call.
//
func (t * Api) parseOrdersJson(body string) ([]types.Order, error) {
  
  var orders []types.Order
  var t_orders []types.TradierOrder  
  
  // Make sure we have at least one account (this should never happen)
  vo := gjson.Get(body, "accounts.account")
  
  if ! vo.Exists() {
    // Return happy (just no orders found)
    return orders, nil	
  }

  // Do we have only one account?
  vo = gjson.Get(body, "accounts.account.account_number")
  
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
  
  // Convert to an formal order array
  t.tempOrderArray2OrderArray(&t_orders, &orders)  
  
  // Return happy
  return orders, nil  
  
}

//
// Parse the case where the user just has more than one account.
// This function just sets the t_orders
//
func (t * Api) orderParseMoreThanOneAccount(body string, t_orders *[]types.TradierOrder) error {
    
  vo := gjson.Get(body, "accounts.account")
  
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
func (t * Api) orderParseOneAccount(body string, t_orders *[]types.TradierOrder) error {
  
  // Do we have any orders.
  vo2 := gjson.Get(body, "accounts.account.orders")
  
  if ! vo2.Exists() {
    return errors.New("No Orders Found") 
  }
  
  // Set the account id
  account_number := gjson.Get(body, "accounts.account.account_number").String() 
  
  // Do we have more than one order
  vo2 = gjson.Get(body, "accounts.account.orders.order.id")
  
  // More than one order??
  if ! vo2.Exists() {
            
    var ws []types.TradierOrder
    
    // Get just the orders part
    vo3 := gjson.Get(body, "accounts.account.orders.order")
    
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
    vo3 := gjson.Get(body, "accounts.account.orders.order")
  
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

//
// Process order data for GetAllOrders
//
func (t * Api) processOrderDataForGetAllOrders(accountNumber string, orders *[]types.Order) (error) {
  
  var t_orders []types.TradierOrder
  
  // Get the JSON
  jsonRt, err := t.SendGetRequest("/accounts/" + accountNumber + "/orders?filter=all")

  if err != nil {
    return err  
  } 
  
  // Do we have more than one order
  vo1 := gjson.Get(jsonRt, "orders.order.Id")
  
  // More than one order??
  if ! vo1.Exists() {
            
    var ws []types.TradierOrder
    
    // Get just the orders part
    vo2 := gjson.Get(jsonRt, "orders.order")
    
    // Unmarshal json
    if err := json.Unmarshal([]byte(vo2.String()), &ws); err != nil {
      return err
    }
    
    // Set the orders to our return
    for _, row := range ws {
      row.AccountId = accountNumber
      t_orders = append(t_orders, row)
    }
  
  } else 
  {    
    var ws types.TradierOrder
  
    // Get just the orders part
    vo2 := gjson.Get(jsonRt, "orders.order")
  
    // Unmarshal json
    if err := json.Unmarshal([]byte(vo2.String()), &ws); err != nil {
      return err
    }
          
    // Set the orders we return.
    ws.AccountId = accountNumber
    t_orders = append(t_orders, ws)   
    
  }  

  // Convert to an formal order array
  t.tempOrderArray2OrderArray(&t_orders, orders)  
  
  // Return the data.
  return nil  
  
}

//
// Go from a temp order array to a formal order array.
//
func (t * Api) tempOrderArray2OrderArray(t_orders *[]types.TradierOrder, orders *[]types.Order) error {
  
  // Clean up the array and make it more generic
  for _, row := range *t_orders {
    
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
    *orders = append(*orders, types.Order{
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
  
  // Success
  return nil  
  
}

/* End File */