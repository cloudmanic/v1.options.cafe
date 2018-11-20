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
  selector: 'app-trading-positions-types-other',
  templateUrl: './other.component.html'
})

export class OtherComponent implements OnInit 
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
  // Get trade group lot count
  //
  getTradeGroupLotCount(tradeGroup: TradeGroup): number 
  {
    if (typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined") 
    {
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
}

/* End File */