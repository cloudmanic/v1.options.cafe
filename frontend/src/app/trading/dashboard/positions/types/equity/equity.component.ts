//
// Date: 9/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Order } from '../../../../../models/order';
import { TradeGroup } from '../../../../../models/trade-group';
import { Position } from '../../../../../models/position';
import { Settings } from '../../../../../models/settings';
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

  @Input() title: string = "";
  @Input() quotes = {};
  @Input() orders: Order[];
  @Input() settings: Settings; 
  @Input() tradeGroups: TradeGroup[];
  actions: DropdownAction[] = null; 

  //
  // Constructor....
  //
  constructor(private tradeService: TradeService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Setup Dropdown actions
    this.setupActions();
  }

  //
  // Setup actions.
  //
  setupActions() 
  {
    let das = []

    // First action
    let da1 = new DropdownAction();
    da1.title = "Close Trade";

    // Place trade to close
    da1.click = (row: TradeGroup) => {

      // Set values
      let tradeDetails = new TradeDetails();
      tradeDetails.Symbol = row.Positions[0].Symbol.ShortName;
      tradeDetails.Class = "equity";
      tradeDetails.Side = "sell";
      tradeDetails.OrderType = "market";
      tradeDetails.Duration = "gtc";
      tradeDetails.Qty = Math.abs(row.Positions[0].Qty);

      // Build legs
      tradeDetails.Legs = [];

      // Open builder to place trade.
      this.tradeService.tradeEvent.emit(new TradeEvent().createNew("toggle-trade-builder", tradeDetails));
    };

    // Add to actions.
    das.push(da1);

    // Load actions.
    this.actions = das;
  }


  //
  // Get market value
  //
  getMarketValue(pos: Position): number
  {
    if (!this.quotes[pos.Symbol.ShortName]) {
      return 0.00
    }

    return this.quotes[pos.Symbol.ShortName].last * pos.Qty    
  }

  //
  // Get total market value.
  //
  getTotalMarketValue(): number
  {
    let total: number = 0.00;

    for (let i = 0; i < this.tradeGroups.length; i++) 
    {
      total = total + this.getMarketValue(this.tradeGroups[i].Positions[0]);
    }

    return total;     
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

    for (let i = 0; i < this.tradeGroups.length; i++) 
    {
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
