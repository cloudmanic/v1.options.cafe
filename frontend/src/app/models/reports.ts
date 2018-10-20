//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';

// Summary Yearly
export class SummaryYearly 
{
  public Year: number;
  public TotalTrades: number;
  public LossCount: number;
  public WinCount: number;
  public Profit: number;
  public Commission: number;
  public WinPercent: number;
  public LossPercent: number;
  public ProfitStd: number;
  public PercentGainStd: number;
  public SharpeRatio: number;
  public AvgRisked: number;
  public AvgPercentGain: number;

  //
  // Json to Object.
  //
  fromJson(json: Object): SummaryYearly {
    let obj = new SummaryYearly();

    obj.Year = json["year"];
    obj.TotalTrades = json["total_trades"];
    obj.LossCount = json["loss_count"];
    obj.WinCount = json["win_count"];
    obj.Profit = json["profit"];
    obj.Commission = json["commission"];
    obj.WinPercent = json["win_percent"];
    obj.LossPercent = json["loss_percent"];
    obj.ProfitStd = json["profit_std"];
    obj.PercentGainStd = json["precent_gain_std"];
    obj.SharpeRatio = json["sharpe_ratio"];
    obj.AvgRisked = json["avg_risked"];
    obj.AvgPercentGain = json["avg_percent_gain"];

    return obj;
  }

  
}

// Profit and Loss
export class ProfitLoss 
{
  Date: Date;
  Profit: number;
  TradeCount: number;
  Commissions: number;
  ProfitPerTrade: number;
  WinRatio: number;
  LossCount: number;
  WinCount: number;

  //
  // Build from json list.
  //
  fromJsonList(json: Object[]): ProfitLoss[] {
    let list: ProfitLoss[] = [];

    if (!json) {
      return list;
    }

    for (let i = 0; i < json.length; i++) {
      list.push(this.fromJson(json[i]));
    }

    return list;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): ProfitLoss 
  {
    let obj = new ProfitLoss();

    obj.Date = moment(json["date"]).toDate();
    obj.Profit = json["profit"];
    obj.TradeCount = json["trade_count"];
    obj.Commissions = json["commissions"];
    obj.ProfitPerTrade = json["profit_per_trade"];
    obj.WinRatio = json["win_ratio"];
    obj.LossCount = json["loss_count"];
    obj.WinCount = json["win_count"];

    return obj;
  }


}

/* End File */