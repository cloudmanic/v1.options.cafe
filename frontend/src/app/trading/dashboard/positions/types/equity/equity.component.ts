//
// Date: 9/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Order } from '../../../../../models/order';
import { TradeGroup } from '../../../../../models/trade-group';
import { Position } from '../../../../../models/position';
import { WebsocketService } from '../../../../../providers/http/websocket.service';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../../../../../providers/http/trade.service';
import { DropdownAction } from '../../../../../shared/dropdown-select/dropdown-select.component';
import { Component, OnInit, Input, Output, OnChanges, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-trading-position-type-equity',
  templateUrl: './equity.component.html',
  styleUrls: []
})

export class EquityComponent implements OnInit {

  @Input() quotes = {};
  @Input() orders: Order[];
  @Input() tradeGroups: TradeGroup[]; 

  //
  // Constructor....
  //
  constructor(private websocketService: WebsocketService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // // Subscribe to data updates from the quotes - Market Quotes
    // this.websocketService.quotePushData.subscribe(data => {
    //   this.quotes[data.symbol] = data;
    // }); 
  }


  //
  // Get average price paid.
  //
  getAvgPricePaid(pos: Position): number 
  {
    return pos.CostBasis / pos.Qty;
  }

  //
  // Total account percent gain.
  //
  getTotalPercentGain(): number 
  {
    return (((this.getTotalGains() + this.getTotalCostBasis()) - this.getTotalCostBasis()) / this.getTotalCostBasis()) * 100
  }

  //
  // Get total daily gain.
  //
  getTotalDailyGain() : number 
  {
    let total: number = 0.00;

    for (let i = 0; i < this.tradeGroups.length; i++) {
      total = total + this.getDailyGain(this.tradeGroups[i].Positions[0]);
    }

    return total;    
  }

  //
  // Get daily Gain
  //
  getDailyGain(pos: Position): number 
  {
    if (!this.quotes[pos.Symbol.ShortName]) {
      return 0.00
    }

    return this.quotes[pos.Symbol.ShortName].change * pos.Qty
  }

  // 
  // Get total cost basis
  //
  getTotalCostBasis() : number
  {
    let total: number = 0.00;

    for(let i = 0; i < this.tradeGroups.length; i++)
    {
      total = total + this.tradeGroups[i].Positions[0].CostBasis;
    }

    return total;
  }

  // 
  // Get total gains
  //
  getTotalGains(): number {
    let total: number = 0.00;

    for (let i = 0; i < this.tradeGroups.length; i++) {
      total = total + this.getTotalGainOfPos(this.tradeGroups[i].Positions[0]);
    }

    return total;
  }

  //
  // Get total gain.
  //
  getTotalGainOfPos(pos: Position): number 
  {
    if(! this.quotes[pos.Symbol.ShortName])
    {
      return 0.00
    }

    return (this.quotes[pos.Symbol.ShortName].last * pos.Qty) - pos.CostBasis
  }  
}

/* End File */
