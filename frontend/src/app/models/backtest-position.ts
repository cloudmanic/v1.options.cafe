//
// Date: 4/12/2022
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Symbol } from './symbol';

//
// BacktestPosition
//
export class BacktestPosition 
{	
  Id: number = 0;
  UserId: number = 0;
  BacktestTradeGroupId: number = 0;
  Status: string = "";
  Symbol: Symbol = new Symbol();
  Qty: number = 0;
  OrgQty: number = 0;
  CostBasis: number = 0.00;
  Proceeds: number = 0.00;
  Profit: number = 0.00;
  AvgOpenPrice: number = 0.00;
  AvgClosePrice: number = 0.00;
  Note: string = "";
  OpenDate: Date = new Date();
  ClosedDate: Date = new Date();

  //
  // Build from JSON list.
  //
  fromJsonList(json: Object[]): BacktestPosition[] 
  {
    let result = [];

    if (!json) 
    {
      return result;
    }

    for (let i = 0; i < json.length; i++) 
    {
      result.push(new BacktestPosition().fromJson(json[i]));
    }

    // Return happy
    return result;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): BacktestPosition 
  {
    let obj = new BacktestPosition();
    
    obj.Id = json["id"];
    obj.BacktestTradeGroupId = json["backtest_trade_group_id"];
    obj.Status = json["status"];
    obj.Symbol = new Symbol().fromJson(json['symbol']);
    obj.Qty = json["qty"];
    obj.OrgQty = json["org_qty"];
    obj.CostBasis = json["cost_basis"];
    obj.Proceeds = json["proceeds"];
    obj.Profit = json["profit"];
    obj.AvgOpenPrice = json["avg_open_price"];
    obj.AvgClosePrice = json["avg_close_price"];
    obj.Note = json["note"];
    obj.OpenDate = moment(json["open_date"]).toDate();
    obj.ClosedDate = moment(json["close_date"]).toDate();

    return obj;
  }  
}

/* End File */