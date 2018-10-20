//
// Date: 10/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import * as moment from 'moment';
import { Settings } from '../../models/settings';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class SettingsService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get settings
  //
  get(): Observable<Settings> 
  {
    return this.http.get<Settings>(environment.app_server + '/api/v1/settings')
      .map((data) => { return new Settings().fromJson(data); });
  }

  //
  // Update Settings
  //
  update(obj: Settings): Observable<boolean>
  {
    let post = {
      strategy_pcs_close_price: obj.StrategyPcsClosePrice,
      strategy_pcs_open_price: obj.StrategyPcsOpenPrice,
      strategy_pcs_lots: obj.StrategyPcsLots,
      strategy_ccs_close_price: obj.StrategyCcsClosePrice,
      strategy_ccs_open_price: obj.StrategyCcsOpenPrice,
      strategy_ccs_lots: obj.StrategyCcsLots,
      strategy_pds_close_price: obj.StrategyPdsClosePrice,
      strategy_pds_open_price: obj.StrategyPdsOpenPrice,
      strategy_pds_lots: obj.StrategyPdsLots,
      strategy_cds_close_price: obj.StrategyCdsClosePrice,
      strategy_cds_open_price: obj.StrategyCdsOpenPrice,
      strategy_cds_lots: obj.StrategyCdsLots,
      notice_trade_filled_email: obj.NoticeTradeFilledEmail,
      notice_trade_filled_sms: obj.NoticeTradeFilledSms,
      notice_trade_filled_push: obj.NoticeTradeFilledPush,
      notice_market_open_email: obj.NoticeMarketOpenedEmail,
      notice_market_open_sms: obj.NoticeMarketOpenedSms,
      notice_market_open_push: obj.NoticeMarketOpenedPush,
      notice_market_closed_email: obj.NoticeMarketClosedEmail,
      notice_market_closed_sms: obj.NoticeMarketClosedSms,
      notice_market_closed_push: obj.NoticeMarketClosedPush
    }

    return this.http.put<boolean>(environment.app_server + '/api/v1/settings', post)
      .map((data) => { return true; });    
  }             
}

/* End File */