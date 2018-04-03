//
// Date: 4/2/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

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
  getOptionExpirations(symbol: string): Observable<Date[]> {
    return this.http.get<string[]>(environment.app_server + '/api/v1/quotes/options/expirations/' + symbol).map(
      (data) => {

        let dates: Date[] = []

        // Build data
        for (let i = 0; i < data.length; i++) {
          dates.push(new Date(data[i]));
        }

        return dates;
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
  Symbol: string;
  Side: string; // buy_to_open, sell_to_open, buy_to_close, sell_to_close
  Qty: number;

  //
  // Create new.
  //
  createNew(symbol: string, side: string, qty: number) : TradeOptionLegs {
    let obj = new TradeOptionLegs();
    obj.Symbol = symbol;
    obj.Side = side;
    obj.Qty = qty;
    return obj;
  }
}

/* End File */