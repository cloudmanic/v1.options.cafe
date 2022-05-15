//
// Date: 4/12/2022
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { BacktestTradeGroup } from './backtest-trade-groups';
import { Screener } from './screener';

//
// Backtest
//
export class Backtest 
{
  Id: number = 0;
  UserId: number = 0;
  Name: string = "";
  StartDate: Date = new Date();
  EndDate: Date = new Date();
  EndingBalance: number = 0.00;
  StartingBalance: number = 0.00;
  CAGR: number = 0.00;
  WinRatio: number = 0.00;
  Return: number = 0.00;
  Profit: number = 0.00;
  TradeCount: number = 0;
  TradeSelect: string = '';
  Midpoint: boolean = false;
  PositionSize: string = '';
  Benchmark: string = '';
  BenchmarkStart: number = 0.00;
  BenchmarkEnd: number = 0.00;
  BenchmarkCAGR: number = 0.00;
  BenchmarkPercent: number = 0.00;
  Screen: Screener = new Screener();
	TradeGroups: BacktestTradeGroup[] = [];
	
  //
  // Build from JSON list.
  //
  fromJsonList(json: Object[]): Backtest[] 
  {
    let result = [];

    if (!json) 
    {
      return result;
    }

    for (let i = 0; i < json.length; i++) 
    {
      result.push(new Backtest().fromJson(json[i]));
    }

    // Return happy
    return result;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): Backtest 
  {
    let obj = new Backtest();

    obj.Id = json["id"];
    obj.UserId = json["user_id"];
    obj.Name = json["name"];
    obj.StartDate = moment(json["start_date"]).toDate();
    obj.EndDate = moment(json["end_date"]).toDate();
    obj.EndingBalance = json["ending_balance"];
    obj.StartingBalance = json["starting_balance"];
    obj.CAGR = json["cagr"];
    obj.WinRatio = json["win_ratio"];
    obj.Return = json["return"];
    obj.Profit = json["profit"];
    obj.TradeCount = json["trade_count"];
    obj.TradeSelect = json["trade_select"];
    obj.Midpoint = json["midpoint"];
    obj.PositionSize = json["position_size"];
    obj.Benchmark = json["benchmark"];
    obj.BenchmarkStart = json["benchmark_start"];
    obj.BenchmarkEnd = json["benchmark_end"];
    obj.BenchmarkCAGR = json["benchmark_cagr"];
    obj.BenchmarkPercent = json["benchmark_percent"];
    obj.Screen = new Screener().fromJson(json['screen']);
    obj.TradeGroups = new BacktestTradeGroup().fromJsonList(json['trade_groups']);

    return obj;
  }  
}

/* End File */