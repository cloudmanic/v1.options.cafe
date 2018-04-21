//
// Date: 4/2/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Symbol } from '../../models/symbol';
import { OptionsChain } from '../../models/options-chain';
import { Observable } from 'rxjs/Observable';
import { HttpClient } from '@angular/common/http';
import { EventEmitter, Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';

@Injectable()
export class TradeService 
{
  tradeEvent = new EventEmitter<TradeEvent>();

  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get option expirations
  //
  getOptionExpirations(symbol: string): Observable<Date[]> 
  {
    return this.http.get<string[]>(environment.app_server + '/api/v1/quotes/options/expirations/' + symbol).map(
      (data) => {

        let dates: Date[] = []

        // Build data
        for (let i = 0; i < data.length; i++) {
          dates.push(moment(data[i]).toDate());
        }

        return dates;
      }
    );
  }

  //
  // Get option expirations
  //
  getOptionStrikesBySymbolExpiration(symbol: string, expire: Date): Observable<number[]> 
  {
    let expr = moment(new Date(expire)).format("YYYY-MM-DD"); 

    return this.http.get<number[]>(environment.app_server + '/api/v1/quotes/options/strikes/' + symbol + '/' + expr).map(
      (data) => {

        let strikes: number[] = []

        // Build data
        for (let i = 0; i < data.length; i++) {
          strikes.push(data[i]);
        }

        return strikes;
      }
    );
  }

  //
  // Push an event.
  //
  PushEvent(event: TradeEvent)
  {
    this.tradeEvent.emit(event);
  }
}

//
// Trade Event
//
export class TradeEvent 
{
  Action: string;
  TradeDetails: TradeDetails;

  //
  // Create new.
  //
  createNew(action: string, tradeDetails: TradeDetails) : TradeEvent
  {
    let obj = new TradeEvent();
    obj.Action = action;
    obj.TradeDetails = tradeDetails;
    return obj;
  }
}

//
// Trade Details
//
export class TradeDetails
{
  Class: string; // equity, option, multileg, combo
  Symbol: string;
  OrderType: string; // market, debit, credit, even
  Duration: string; // day, gtc
  Price: number;
  Legs: TradeOptionLegs[];
}

//
// Trade Option Legs
//
export class TradeOptionLegs 
{
  Symbol: Symbol;
  Expire: Date; // Used just for the forms (not to be posted when creating a trade)
  Strike: number;
  Side: string; // buy_to_open, sell_to_open, buy_to_close, sell_to_close
  Qty: number;
  Type: string; // (Put | Call) Used just for the forms (not to be posted when creating a trade)
  Strikes: number[]; // Used just for the forms (not to be posted when creating a trade)
  Chain: OptionsChain; // Used just for the forms (not to be posted when creating a trade)

  //
  // Create new.
  //
  createNew(symbol: Symbol, expire: Date, type: string, strike: number, side: string, qty: number): TradeOptionLegs {
    let obj = new TradeOptionLegs();
    obj.Symbol = symbol;
    obj.Expire = expire;
    obj.Type = type;
    obj.Strike = strike;
    obj.Side = side;
    obj.Qty = qty;
    obj.Strikes = [];
    return obj;
  }
}

/* End File */