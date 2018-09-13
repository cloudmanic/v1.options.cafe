//
// Date: 9/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';

export class BrokerEvent 
{
  public Id: number;
  public Date: Date;
  public Type: string;
  public Amount: number;
  public Symbol: string;
  public Commission: number;
  public Description: string;
  public Price: number;
  public Quantity: number;
  public TradeType: string;

  //
  // Build from JSON list.
  //
  fromJsonList(json: Object[]): BrokerEvent[] {
    let result = [];

    if (!json) 
    {
      return result;
    }

    for (let i = 0; i < json.length; i++) 
    {
      result.push(new BrokerEvent().fromJson(json[i]));
    }

    // Return happy
    return result;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): BrokerEvent {
    let obj = new BrokerEvent();

    obj.Id = json["id"];
    obj.Date = moment(json["date"]).toDate();
    obj.Type = json["type"];
    obj.Amount = json["amount"];
    obj.Symbol = json["symbol"];
    obj.Commission = json["commission"];
    obj.Price = json["price"];
    obj.Quantity = json["quantity"];
    obj.TradeType = json["trade_type"];
    obj.Description = json["description"];

    return obj;
  }

  
}

/* End File */