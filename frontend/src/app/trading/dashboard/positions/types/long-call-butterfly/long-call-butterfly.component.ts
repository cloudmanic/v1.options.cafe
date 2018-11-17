//
// Date: 11/15/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input, Output, OnChanges, EventEmitter } from '@angular/core';
import { Order } from '../../../../../models/order';
import { TradeGroup } from '../../../../../models/trade-group';
import { Position } from '../../../../../models/position';
import { Settings } from '../../../../../models/settings';
import { DropdownAction } from '../../../../../shared/dropdown-select/dropdown-select.component';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../../../../../providers/http/trade.service';
import { AnalyzeService, AnalyzeTrade } from '../../../../../providers/http/analyze.service';
import { AnalyzeLeg } from '../../../../../models/analyze-result';

@Component({
  selector: 'app-trading-positions-long-call-butterfly',
  templateUrl: './long-call-butterfly.component.html',
  styleUrls: []
})

export class LongCallButterflyComponent implements OnInit 
{
  @Input() quotes = {};
  @Input() orders: Order[];
  @Input() title: string;
  @Input() settings: Settings;
  @Input() tradeGroups: TradeGroup[];

  actions: DropdownAction[] = null;

  //
  // Constructor....
  //
  constructor(private tradeService: TradeService, private analyzeService: AnalyzeService) { }

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
    // Build social section
    let closeSection = new DropdownAction();
    closeSection.title = "Close Position";
    closeSection.section = true;

    // First action
    let da1 = new DropdownAction();
    da1.title = "Close Trade";

    // Place trade to close
    da1.click = (row: TradeGroup) => {

      // Set values
      let tradeDetails = new TradeDetails();
      tradeDetails.Symbol = row.Positions[0].Symbol.OptionUnderlying;
      tradeDetails.Class = "multileg";
      tradeDetails.OrderType = "debit";
      tradeDetails.Duration = "gtc";

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


    // Close Trade @ Market
    let da3 = new DropdownAction();
    da3.title = "Close @ Market";

    // Place trade to close
    da3.click = (row: TradeGroup) => {

      // Set values
      let tradeDetails = new TradeDetails();
      tradeDetails.Symbol = row.Positions[0].Symbol.OptionUnderlying;
      tradeDetails.Class = "multileg";
      tradeDetails.OrderType = "market";
      tradeDetails.Duration = "gtc";

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

    // Build Review section
    let reviewSection = new DropdownAction();
    reviewSection.title = "Review Position";
    reviewSection.section = true;

    // Analyze Trade
    let analyze = new DropdownAction();
    analyze.title = "Analyze Trade";

    // Analyze Trade Click
    analyze.click = (row: TradeGroup) => {
      this.analyzeTrade(row)
    };

    // Build social section
    let socialSection = new DropdownAction();
    socialSection.title = "Share Position";
    socialSection.section = true;

    // Tweet Trade
    let tweet = new DropdownAction();
    tweet.title = "Tweet Trade";

    // Place trade to close
    tweet.click = (row: TradeGroup) => {

      let debit = (row.Risked / 100) / row.Positions[0].Qty;

      let tweet = "I just opened a new " + row.Type + " today on the " + row.Positions[0].Symbol.OptionUnderlying + ". For a debit of $" + debit + ".%0a%0a";

      for (let i = 0; i < row.Positions.length; i++) {
        tweet = tweet + row.Positions[i].Symbol.Name + "%0a";
      }

      tweet = tweet + "%0a";

      window.open('https://twitter.com/share?text=' + tweet + '&via=options_cafe&url=https://options.cafe&hashtags=OptionsTrading', '', 'menubar=no, toolbar = no, resizable = yes, scrollbars = yes, height = 600, width = 600');
    };

    // Load actions.
    this.actions = [closeSection, da1, da3, reviewSection, analyze, socialSection, tweet];
  }

  //
  // Analyze Trade
  //
  analyzeTrade(tradeGroup: TradeGroup)
  {
    let trade = new AnalyzeTrade();
    trade.OpenCost = tradeGroup.Risked;
    trade.Legs = [];

    // Add Legs
    for(let i = 0; i < tradeGroup.Positions.length; i++)
    {
      let leg = new AnalyzeLeg();
      leg.Qty = tradeGroup.Positions[i].Qty;
      leg.SymbolStr = tradeGroup.Positions[i].Symbol.ShortName;
      trade.Legs.push(leg) 
    }

    // Send request to show analyze dialog
    this.analyzeService.dialog.emit(trade);
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
  // Get the trade group widget total.
  //
  getTradeGroupWidgetTotal(tradeGroups: TradeGroup[]): number 
  {
    let total = 0.00;

    for(var i = 0; i < tradeGroups.length; i++) 
    {
      total = total + tradeGroups[i].Risked;
    }

    return total;
  }

  //
  // Get trade group lot count
  //
  getTradeGroupLotCount(tradeGroup: TradeGroup): number 
  {
    if(typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined") 
    {
      return 0;
    }

    return Math.abs(tradeGroup.Positions[0].Qty);
  }

  //
  // Get trade profit and loss
  //
  getTradeProfitLoss(tradeGroup: TradeGroup): number 
  {
    let close = 0.00;

    for (var i = 0; i < tradeGroup.Positions.length; i++) 
    {
      if (typeof this.quotes[tradeGroup.Positions[i].Symbol.ShortName] == "undefined") 
      {
        return 0;
      }

      if (typeof this.quotes[tradeGroup.Positions[i].Symbol.ShortName].bid == "undefined") 
      {
        return 0;
      }

      if (typeof this.quotes[tradeGroup.Positions[i].Symbol.ShortName].ask == "undefined") 
      {
        return 0;
      }

      if (tradeGroup.Positions[i].Qty > 0)
      {
        close = close + (this.quotes[tradeGroup.Positions[i].Symbol.ShortName].bid * tradeGroup.Positions[i].Qty * 100);
      } else 
      {
        close = close + (this.quotes[tradeGroup.Positions[i].Symbol.ShortName].ask * tradeGroup.Positions[i].Qty * 100);
      }
    }    

    return close - tradeGroup.Risked;
  }
}

/* End File */
