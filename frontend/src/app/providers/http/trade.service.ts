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
  // Cancel Order
  //
  cancelOrder(brokerAccountId: number, brokerOrderId: number): Observable<boolean>
  {
    return this.http.delete<boolean>(environment.app_server + '/api/v1/orders/' + brokerAccountId + '/' + brokerOrderId, {})
      .map((data) => { return true; });    
  }

  //
  // Submit trade.
  //
  submitTrade(trade: TradeDetails, brokerAccountId: string): Observable<OrderSubmit>
  {
    let body = {
      broker_account_id: parseInt(brokerAccountId),
      side: trade.Side,
      quantity: Number(trade.Qty),      
      class: trade.Class,
      symbol: trade.Symbol,
      duration: trade.Duration,
      type: trade.OrderType,
      price: Number(trade.Price),
      stop: Number(trade.Stop),
      legs: []
    }

    if(trade.Legs) 
    {
      for (let i = 0; i < trade.Legs.length; i++) 
      {
        body.legs.push(new TradeOptionLegsPost().createNew(
          trade.Legs[i].Side,
          Number(trade.Legs[i].Qty),
          trade.Legs[i].Symbol.ShortName
        ));
      }
    }

    return this.http.post<OrderSubmit>(environment.app_server + '/api/v1/orders', body)
      .map((data) => { return new OrderSubmit().fromJson(data); });
  }

  //
  // Preview trade.
  //
  previewTrade(trade: TradeDetails, brokerAccountId: string): Observable<OrderPreview>
  {
    let body = {
      broker_account_id: parseInt(brokerAccountId),
      side: trade.Side,
      quantity: Number(trade.Qty),
      class: trade.Class,
      symbol: trade.Symbol,
      duration: trade.Duration,
      type: trade.OrderType,
      price: Number(trade.Price),
      stop: Number(trade.Stop),
      legs: []
    }

    if(trade.Legs) 
    {
      for (let i = 0; i < trade.Legs.length; i++) 
      {
        body.legs.push(new TradeOptionLegsPost().createNew(
          trade.Legs[i].Side,
          Number(trade.Legs[i].Qty),
          trade.Legs[i].Symbol.ShortName
        ));
      }
    }

    return this.http.post<OrderPreview>(environment.app_server + '/api/v1/orders/preview', body)
      .map((data) => { return new OrderPreview().fromJson(data); });   
  }

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
// Trade Submit response
//
export class OrderSubmit {
  Id: number;
  Status: string;
  Error: string;

  //
  // Json to Object.
  //
  fromJson(json: Object): OrderSubmit {
    let op = new OrderSubmit();
    op.Id = json["id"];
    op.Status = json["status"];
    op.Error = json["error"];
    return op;
  }
}

//
// Trade Preview response
//
export class OrderPreview 
{
  Status: string;
  Error: string;
  Commission: number;
  Cost: number;
  Fees: number;
  Symbol: string;
  Type: string;
  Duration: string;
  Price: number;
  OrderCost: number;
  MarginChange: number;
  OptionRequirement: number;
  Class: string;
  Strategy: string;

  //
  // Json to Object.
  //
  fromJson(json: Object): OrderPreview {
    let op = new OrderPreview();
    op.Status = json["status"];
    op.Error = json["error"];
    op.Commission = json["commission"];
    op.Cost = json["cost"];
    op.Fees = json["fees"];
    op.Symbol = json["symbol"];
    op.Type = json["type"];
    op.Duration = json["duration"];
    op.Price = json["price"];
    op.OrderCost = json["order_cost"];
    op.MarginChange = json["margin_change"];
    op.OptionRequirement = json["option_requirement"];
    op.Class = json["class"];
    op.Strategy = json["strategy"];
    return op;
  }   
}

//
// Trade Details
//
export class TradeDetails
{
  Class: string; // equity, option, multileg, combo
  Symbol: string;
  Side: string; // buy, sell, buy_to_cover, sell_short
  OrderType: string; // market, debit, credit, even, limit, stop, stop_limit
  Duration: string; // day, gtc
  Price: number;
  Stop: number;
  Qty: number;
  Legs: TradeOptionLegs[];
}

//
// Trade Option Legs
//
export class TradeOptionLegsPost 
{
  option_symbol: string;
  side: string;
  quantity: number;

  //
  // Create new.
  //
  createNew(side: string, qty: number, symbol: string): TradeOptionLegsPost {
    let obj = new TradeOptionLegsPost();
    obj.option_symbol = symbol;
    obj.side = side;
    obj.quantity = qty;
    return obj;
  }  
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