//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { OrderLeg } from './order-leg';

export class Order {
  constructor(
    public Id: number,
    public AccountId: string,
    public AvgFillPrice: number,
    public Class: string,
    public CreateDate: string,
    public Duration: string,
    public ExecQuantity: string,
    public LastFillPrice: number,
    public LastFillQuantity: number,
    public NumLegs: number,
    public Price: number,
    public Quantity: number,
    public RemainingQuantity: number,
    public Side: string,
    public Status: string,
    public Symbol: string,
    public TransactionDate: string,
    public Type: string,
    public OptionSymbol: string,
    public Legs: OrderLeg[]
  ){}

  //
  // Build build the data for emitting to the app. 
  //
  public static buildForEmit(data) : Order[] {
    
    let orders = [];
    
    // Sometimes we do not have data yet.
    if(! data)
    {
      return orders;
    }
    
    for(let i = 0; i < data.length; i++)
    {
      // Add in the legs
      let legs = [];
      
      if(data[i].legs)
      {
        for(let k = 0; k < data[i].legs.length; k++)
        {
          legs.push(new OrderLeg(
            data[i].legs[k].class,
            data[i].legs[k].symbol,
            data[i].legs[k].option_symbol, 
            data[i].legs[k].side, 
            data[i].legs[k].quantity, 
            data[i].legs[k].status, 
            data[i].legs[k].duration, 
            data[i].legs[k].avg_fill_price, 
            data[i].legs[k].exec_quantity, 
            data[i].legs[k].last_fill_price, 
            data[i].legs[k].last_fill_quantity, 
            data[i].legs[k].remaining_quantity, 
            data[i].legs[k].create_date, 
            data[i].legs[k].transaction_date          
          ));
        }
      }
      
      // Push the order on
      orders.push(new Order(
          data[i].id,
          data[i].account_id,
          data[i].avg_fill_price,
          data[i].class,
          data[i].create_date,
          data[i].duration,
          data[i].exec_quantity,
          data[i].last_fill_price,
          data[i].last_fill_quantity,
          data[i].num_legs,
          data[i].price,
          data[i].quantity,
          data[i].remaining_quantity,
          data[i].side,
          data[i].status,
          data[i].symbol,
          data[i].transaction_date,
          data[i].type,
          data[i].option_symbol,
          legs));
               
    }
    
    return orders;
        
  } 

}

/* End File */