//
// Date: 10/11/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

export class Settings 
{
  Id: number;

  // Put credit Spread
  StrategyPcsClosePrice: number;
  StrategyPcsOpenPrice: string;
  StrategyPcsLots: number;

  // Call credit Spread
  StrategyCcsClosePrice: number;
  StrategyCcsOpenPrice: string;
  StrategyCcsLots: number;

  // Put debit Spread
  StrategyPdsClosePrice: number;
  StrategyPdsOpenPrice: string;
  StrategyPdsLots: number;

  // Call debit Spread
  StrategyCdsClosePrice: number;
  StrategyCdsOpenPrice: string;
  StrategyCdsLots: number;

  // Trade Filled
  NoticeTradeFilledEmail: string;
  NoticeTradeFilledSms: string;
  NoticeTradeFilledPush: string;  

  // Market Open 
  NoticeMarketOpenedEmail: string;
  NoticeMarketOpenedSms: string;
  NoticeMarketOpenedPush: string;  
 
  // Market Closed 
  NoticeMarketClosedEmail: string;
  NoticeMarketClosedSms: string;
  NoticeMarketClosedPush: string;  

  //
  // Build from JSON.
  //
  fromJson(json: Object): Settings 
  {
    let result = new Settings();

    if (! json) 
    {
      return result;
    }

    // Put credit Spread
    result.StrategyPcsClosePrice = json["strategy_pcs_close_price"];
    result.StrategyPcsOpenPrice = json["strategy_pcs_open_price"];
    result.StrategyPcsLots = json["strategy_pcs_lots"];

    // Call credit Spread  
    result.StrategyCcsClosePrice = json["strategy_ccs_close_price"];
    result.StrategyCcsOpenPrice = json["strategy_ccs_open_price"];
    result.StrategyCcsLots = json["strategy_ccs_lots"];

    // Put debit Spread
    result.StrategyPdsClosePrice = json["strategy_pds_close_price"];
    result.StrategyPdsOpenPrice = json["strategy_pds_open_price"];
    result.StrategyPdsLots = json["strategy_pds_lots"];

    // Call debit Spread
    result.StrategyCdsClosePrice = json["strategy_cds_close_price"];
    result.StrategyCdsOpenPrice = json["strategy_cds_open_price"];
    result.StrategyCdsLots = json["strategy_cds_lots"];

    // Trade Filled
    result.NoticeTradeFilledEmail = json["notice_trade_filled_email"];
    result.NoticeTradeFilledSms = json["notice_trade_filled_sms"];
    result.NoticeTradeFilledPush = json["notice_trade_filled_push"];

    // Market Open
    result.NoticeMarketOpenedEmail = json["notice_market_open_email"];
    result.NoticeMarketOpenedSms = json["notice_market_open_sms"];
    result.NoticeMarketOpenedPush = json["notice_market_open_push"];

    // Market Closed
    result.NoticeMarketClosedEmail = json["notice_market_closed_email"];
    result.NoticeMarketClosedSms = json["notice_market_closed_sms"];
    result.NoticeMarketClosedPush = json["notice_market_closed_push"];  

    // Return happy
    return result;
  }
}



/* End File */
