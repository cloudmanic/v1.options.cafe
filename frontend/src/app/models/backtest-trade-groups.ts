//
// Date: 4/12/2022
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { BacktestPosition } from './backtest-position';

//
// BacktestTradeGroup
//
export class BacktestTradeGroup 
{
  Id: number = 0;
  UserId: number = 0;
  BacktestId: number = 0;
  Strategy: string = "";
  Status: string = "";
  SpreadText: string = "";
  OpenDate: Date = new Date();
  CloseDate: Date = new Date();
  OpenPrice: number = 0.00;
  ClosePrice: number = 0.00;
  Credit: number = 0.00;
  ReturnPercent: number = 0.00;
  ReturnFromStart: number = 0.00;
  Margin: number = 0.00;
  Balance: number = 0.00;
  PutPrecentAway: number = 0.00;
  CallPrecentAway: number = 0.00;
  BenchmarkLast: number = 0.00;
  BenchmarkBalance: number = 0.00;
  BenchmarkReturn: number = 0.00;
  Lots: number = 0;
  Positions: BacktestPosition[] = [];
  Note: string = "";

  //
  // Build from JSON list.
  //
  fromJsonList(json: Object[]): BacktestTradeGroup[] 
  {
    let result = [];

    if (!json) 
    {
      return result;
    }

    for (let i = 0; i < json.length; i++) 
    {
      result.push(new BacktestTradeGroup().fromJson(json[i]));
    }

    // Return happy
    return result;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): BacktestTradeGroup 
  {
    let obj = new BacktestTradeGroup();

    obj.Id = json["id"];
    obj.UserId = json["user_id"];
    obj.BacktestId = json["backtest_id"];
    obj.Strategy = json["strategy"];
    obj.SpreadText = json["spread_text"];
    obj.Status = json["status"];
    obj.OpenDate = moment(json["open_date"]).toDate();
    obj.CloseDate = moment(json["close_date"]).toDate();
    obj.OpenPrice = json["open_price"];
    obj.ClosePrice = json["close_price"];
    obj.Credit = json["credit"];
    obj.ReturnPercent = json["return_percent"];
    obj.ReturnFromStart = json["return_from_start"];
    obj.Margin = json["margin"];
    obj.Balance = json["balance"];
    obj.PutPrecentAway = json["put_percent_away"];
    obj.CallPrecentAway = json["call_percent_away"];
    obj.BenchmarkLast = json["benchmark_last"];
    obj.BenchmarkBalance = json["benchmark_balance"];
    obj.BenchmarkReturn = json["benchmark_return"];
    obj.Lots = json["lots"];
    obj.Note = json["note"];
    obj.Positions = new BacktestPosition().fromJsonList(json['positions']);

    return obj;
  }  
}

/* End File */