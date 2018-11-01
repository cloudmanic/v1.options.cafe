//
// Date: 2/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input } from '@angular/core';
import { DropdownAction } from '../../../../../shared/dropdown-select/dropdown-select.component';
import { Position } from '../../../../../models/position';
import { TradeGroup } from '../../../../../models/trade-group';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../../../../../providers/http/trade.service';
import { Order } from '../../../../../models/order';

@Component({
  selector: 'app-trading-positions-types-option',
  templateUrl: './option.component.html'
})

export class OptionComponent implements OnInit 
{
  @Input() quotes = {};
  @Input() orders: Order[];
  @Input() tradeGroups: TradeGroup[];

  actions: DropdownAction[] = null;

  //
  // Constructor....
  //
  constructor(private tradeService: TradeService) { }

  //
  // OnInit....
  //
  ngOnInit() {

    // Setup Dropdown actions
    this.setupActions();

  }

  //
  // Setup actions.
  //
  setupActions() {
    let das = []

    // First action
    let da1 = new DropdownAction();
    da1.title = "Close Trade @ $0.03";

    // Place trade to close
    da1.click = (row: TradeGroup) => {

      // Set values
      let tradeDetails = new TradeDetails();
      tradeDetails.Symbol = "SPY";
      tradeDetails.Class = "multileg";
      tradeDetails.OrderType = "debit";
      tradeDetails.Duration = "gtc";
      tradeDetails.Price = 0.03;

      // Build legs
      tradeDetails.Legs = [];

      for (let i = 0; i < row.Positions.length; i++) 
      {
        let side = "sell_to_close";
        let qty = row.Positions[i].Qty;

        if (row.Positions[i].Qty < 0) 
        {
          side = "buy_to_close";
          qty = qty * -1;
        }

        tradeDetails.Legs.push(new TradeOptionLegs().createNew(row.Positions[i].Symbol, row.Positions[i].Symbol.OptionExpire, row.Positions[i].Symbol.OptionType, row.Positions[i].Symbol.OptionStrike, side, qty));
      }

      // Open builder to place trade.
      this.tradeService.tradeEvent.emit(new TradeEvent().createNew("toggle-trade-builder", tradeDetails));
    };

    das.push(da1);

    this.actions = das;
  }

  //
  // Get trade group days to expire
  //
  getTradeGroupDaysToExpire(tradeGroup: TradeGroup): number 
  {
    if (typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined") 
    {
      return 0;
    }

    return Math.round((tradeGroup.Positions[0].Symbol.OptionExpire.getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24));
  }

  //
  // Get trade group percent away
  //
  getTradeGroupPercentAway(tradeGroup: TradeGroup): number 
  {
    // Find the short strike.
    var short_strike = null;

    for (let i = 0; i < tradeGroup.Positions.length; i++) 
    {
      if (tradeGroup.Positions[i].Symbol.Type != "Option") {
        continue;
      }

      if (typeof this.quotes[tradeGroup.Positions[i].Symbol.ShortName] == "undefined") {
        return 0.00;
      }

      if (tradeGroup.Positions[i].Qty < 0) {
        short_strike = tradeGroup.Positions[i];
      }
    }

    if (short_strike == null) {
      return 0.00;
    }

    if (tradeGroup.Positions[0].Symbol.OptionType == 'Put') {
      return ((this.quotes[short_strike.Symbol.OptionUnderlying].last - short_strike.Symbol.OptionStrike) /
        ((this.quotes[short_strike.Symbol.OptionUnderlying].last + short_strike.Symbol.OptionStrike) / 2)) * 100;
    } else {
      return ((short_strike.Symbol.OptionStrike - this.quotes[short_strike.Symbol.OptionUnderlying].last) /
        ((short_strike.Symbol.OptionStrike + this.quotes[short_strike.Symbol.OptionUnderlying].last) / 2)) * 100;
    }
  }

  //
  // Get the trade group widget total.
  //
  getTradeGroupWidgetTotal(tradeGroups: TradeGroup[]): number 
  {
    let total: number = 0.00

    // Loop through the tradegroups and add them up.
    for (let i = 0; i < tradeGroups.length; i++) 
    {
      for (let k = 0; k < tradeGroups[i].Positions.length; k++) 
      {
        total += this.getTradeProfitLoss(tradeGroups[i]);
      }
    }

    return total;
  }

  //
  // Get trade group lot count
  //
  getTradeGroupLotCount(tradeGroup: TradeGroup): number {
    if (typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined") {
      return 0;
    }

    return Math.abs(tradeGroup.Positions[0].Qty);
  }

  //
  // Get Single Value (based on a bid or an ask)
  //
  getSingleValue(position: Position): number 
  {
    if (typeof this.quotes[position.Symbol.ShortName] == "undefined") 

    {
      return 0.00;
    }

    // Short or long?
    if (position.Qty > 0) 
    {
      return this.quotes[position.Symbol.ShortName].bid * position.Qty * 100;
    } else 
    {
      return this.quotes[position.Symbol.ShortName].ask * position.Qty * 100;
    }
  }

  //
  // Progress bar for a Option trade
  //
  getProgressOptionsTrade(tradeGroup: TradeGroup): number 
  {
    if (typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined") 
    {
      return 0.00;
    }

    var order: Order = null;

    // Loop through open orders and find this trade
    for (let i = 0; i < this.orders.length; i++) 
    {
      if (this.orders[i].OptionSymbol == tradeGroup.Positions[0].Symbol.ShortName) 
      {
        order = this.orders[i];
      }
    }

    if (order) 
    {
      // Short or long?
      if (tradeGroup.Positions[0].Qty > 0) 
      {
        let open_price = (tradeGroup.Positions[0].CostBasis / tradeGroup.Positions[0].Qty) / 100;
        let top = (this.quotes[tradeGroup.Positions[0].Symbol.ShortName].bid - open_price)
        return (top / order.Price) * 100;
      }
    }

    return 0.00;
  }

  //
  // Get trade progress
  //
  getTradeProgress(tradeGroup: TradeGroup): number 
  {
    return this.getProgressOptionsTrade(tradeGroup);
  }

  //
  // Get trade progress for the bar
  //
  getTradeProgressBar(tradeGroup: TradeGroup): number 
  {
    let p = 0.00;
    p = this.getProgressOptionsTrade(tradeGroup);

    // keep it within a range
    if ((p > 0) && (p <= 100)) 
    {
      return p;
    } else if (p > 100) 
    {
      return 100;
    }

    return 0.00
  }

  //
  // Get the total P&L for put credit spreads
  //
  getPutCreditSpreadProfitLoss(tradeGroup: TradeGroup): number {
    if (typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined") {
      return 0.00;
    }

    if (typeof this.quotes[tradeGroup.Positions[1].Symbol.ShortName] == "undefined") {
      return 0.00;
    }

    return tradeGroup.Credit - (((this.quotes[tradeGroup.Positions[1].Symbol.ShortName].ask - this.quotes[tradeGroup.Positions[0].Symbol.ShortName].bid) * 100) * Math.abs(tradeGroup.Positions[0].Qty));
  }

  //
  // Get the total P&L for put credit spreads - Call Credit Spread
  //
  getCallCreditSpreadProfitLoss(tradeGroup: TradeGroup): number 
  {
    if (typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined") 
    {
      return 0.00;
    }

    if (typeof this.quotes[tradeGroup.Positions[1].Symbol.ShortName] == "undefined") 
    {
      return 0.00;
    }

    return tradeGroup.Credit - ((((this.quotes[tradeGroup.Positions[1].Symbol.ShortName].ask - this.quotes[tradeGroup.Positions[0].Symbol.ShortName].bid) * 100) * Math.abs(tradeGroup.Positions[0].Qty)) * -1);
  }

  //
  // Get trade profit and loss
  //
  getTradeProfitLoss(tradeGroup: TradeGroup): number 
  {
    if (typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined") 
    {
      return 0.00;
    }

    return (this.quotes[tradeGroup.Positions[0].Symbol.ShortName].last * tradeGroup.Positions[0].Qty * 100) - tradeGroup.Positions[0].CostBasis;    
  }

  //
  // Do we show a progress bar
  //
  showProgressbar(): boolean 
  {
    return true;
  }
}

/* End File */