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
    public Legs: OrderLeg[]
  ){}

  //
  // Build build the data for emitting to the app. 
  //
  public static buildForEmit(data) : Order[] {
    
    let orders = [];
    
    for(let i = 0; i < data.length; i++)
    {
      // Add in the legs
      let legs = [];
      
      if(data[i].NumLegs > 0)
      {
        for(let k = 0; k < data[i].Legs.length; k++)
        {
          legs.push(new OrderLeg(
            data[i].Legs[k].Type,
            data[i].Legs[k].Symbol,
            data[i].Legs[k].OptionSymbol, 
            data[i].Legs[k].Side, 
            data[i].Legs[k].Quantity, 
            data[i].Legs[k].Status, 
            data[i].Legs[k].Duration, 
            data[i].Legs[k].AvgFillPrice, 
            data[i].Legs[k].ExecQuantity, 
            data[i].Legs[k].LastFillPrice, 
            data[i].Legs[k].LastFillQuantity, 
            data[i].Legs[k].RemainingQuantity, 
            data[i].Legs[k].CreateDate, 
            data[i].Legs[k].TransactionDate          
          ));
        }
      }
      
      // Push the order on
      orders.push(new Order(
          data[i].Id,
          data[i].AccountId,
          data[i].AvgFillPrice,
          data[i].Class,
          data[i].CreateDate,
          data[i].Duration,
          data[i].ExecQuantity,
          data[i].LastFillPrice,
          data[i].LastFillQuantity,
          data[i].NumLegs,
          data[i].Price,
          data[i].Quantity,
          data[i].RemainingQuantity,
          data[i].Side,
          data[i].Status,
          data[i].Symbol,
          data[i].TransactionDate,
          data[i].Type,
          legs));
               
    }
    
    return orders;
        
  } 

}

/* End File */