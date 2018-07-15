//
// Date: 2/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Order } from '../../../models/order';
import { TradeGroup } from '../../../models/trade-group';
import { Position } from '../../../models/position';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../../../providers/http/trade.service';
import { DropdownAction } from '../../../shared/dropdown-select/dropdown-select.component';
import { Component, OnInit, Input, Output, OnChanges, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-trading-position',
  templateUrl: './position.component.html'
})

export class PositionComponent implements OnInit 
{
  @Input() quotes = {};
  @Input() orders: Order[];  
  @Input() title: string;   
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
  setupActions()
  {
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
  // Get trade group total header title.
  //
  getTradeGroupTotalHeaderTitle(tradeGroups: TradeGroup[]) : string 
  {
    // Get progress based on type.
    switch(this.title)
    {
      case 'Put Credit Spreads':
      case 'Call Credit Spreads':
        return "Credit";

      case 'Options':
        return "P&amp;L";             
    }

    return "";
  }

  //
  // Get trade group days to expire
  //
  getTradeGroupDaysToExpire(tradeGroup: TradeGroup) : number
  {
    if(typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined")
    {
      return 0;
    }   
    
    return Math.round((tradeGroup.Positions[0].Symbol.OptionExpire.getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24));   
  }

  //
  // Get trade group percent away
  //
  getTradeGroupPercentAway(tradeGroup: TradeGroup) : number
  {
    // Find the short strike.
    var short_strike = null; 

    for(let i = 0; i < tradeGroup.Positions.length; i++)
    {
      if(tradeGroup.Positions[i].Symbol.Type != "Option")
      {
        continue;
      }

      if (typeof this.quotes[tradeGroup.Positions[i].Symbol.ShortName] == "undefined") 
      {
        return 0.00;
      }       

      if(tradeGroup.Positions[i].Qty < 0)
      {
        short_strike = tradeGroup.Positions[i];
      }
    }

    if(short_strike == null)
    {
       return 0.00;
    }

    if (tradeGroup.Positions[0].Symbol.OptionType == 'Put')
    {
      return ((this.quotes[short_strike.Symbol.OptionUnderlying].last - short_strike.Symbol.OptionStrike) / 
        ((this.quotes[short_strike.Symbol.OptionUnderlying].last + short_strike.Symbol.OptionStrike) / 2)) * 100;
    } else
    {
      return ((short_strike.Symbol.OptionStrike - this.quotes[short_strike.Symbol.OptionUnderlying].last) / 
        ((short_strike.Symbol.OptionStrike + this.quotes[short_strike.Symbol.OptionUnderlying].last) / 2)) * 100;    
    }
  }

  //
  // Get the trade group widget total.
  //
  getTradeGroupWidgetTotal(tradeGroups: TradeGroup[]) : number
  {
    switch(this.title)
    {
      case 'Put Credit Spreads':
      case 'Call Credit Spreads':
        return this.getCreditSpreadWidgetProfitLoss(tradeGroups);

      case 'Options':
        return this.getOptionWidgetProfitLoss(tradeGroups);
    }

    // Default to blank
    return 0.00;
  }

  //
  // Get the total credit for the credit spread widget
  //
  getCreditSpreadWidgetProfitLoss(tradeGroups: TradeGroup[]) : number 
  {
    let total = 0.00;

    for(var i = 0; i < tradeGroups.length; i++)
    {
      total = total + tradeGroups[i].Credit;
    }

    return total;
  }

  //
  // Get trade group lot count
  //
  getTradeGroupLotCount(tradeGroup: TradeGroup) : number
  {
    if(typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined")
    {
      return 0;
    }

    return Math.abs(tradeGroup.Positions[0].Qty);
  }

  //
  // Get Single Value (based on a bid or an ask)
  //
  getSingleValue(position: Position) : number 
  {
    if(typeof this.quotes[position.Symbol.ShortName] == "undefined")
    {
      return 0.00;
    }

    // Short or long?
    if(position.Qty > 0)
    {
      return this.quotes[position.Symbol.ShortName].bid * position.Qty * 100;
    } else 
    {
      return this.quotes[position.Symbol.ShortName].ask * position.Qty * 100;
    }     
  }

  //
  // Get the total P&L for options
  //
  getOptionWidgetProfitLoss(tradeGroups: TradeGroup[]) : number 
  {
    let total: number = 0.00

    // Loop through the tradegroups and add them up.
    for(let i = 0; i < tradeGroups.length; i++)
    {
      for(let k = 0; k < tradeGroups[i].Positions.length; k++)
      {      
        total += this.getOptionProfitLoss(tradeGroups[i].Positions[k]);
      }
    }

    return total;
  }

  //
  // Figure out the Profit & loss for a position - Option
  //
  getOptionProfitLoss(position: Position) : number 
  {
    if(typeof this.quotes[position.Symbol.ShortName] == "undefined")
    {
      return 0.00;
    }

    return (this.quotes[position.Symbol.ShortName].last * position.Qty * 100) - position.CostBasis;
  }

  //
  // Progress bar for a Option trade
  //
  getProgressOptionsTrade(tradeGroup: TradeGroup) : number 
  {
    if(typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined")
    {
      return 0.00;
    }

    var order: Order = null;

    // Loop through open orders and find this trade
    for(let i = 0; i < this.orders.length; i++)
    {
      if(this.orders[i].OptionSymbol == tradeGroup.Positions[0].Symbol.ShortName)
      {
        order = this.orders[i];
      }
    }

    if(order)
    {
      // Short or long?
      if(tradeGroup.Positions[0].Qty > 0)
      {
        let open_price = (tradeGroup.Positions[0].CostBasis / tradeGroup.Positions[0].Qty) / 100;
        let top = (this.quotes[tradeGroup.Positions[0].Symbol.ShortName].bid - open_price)
        return (top  / order.Price) * 100;
      }
    }

    return 0.00;
  }

  //
  // Get trade progress
  //
  getTradeProgress(tradeGroup: TradeGroup) : number
  {    
    switch(tradeGroup.Type)
    {
      case 'Option':
        return this.getProgressOptionsTrade(tradeGroup);

      case 'Put Credit Spread':
        return (this.getPutCreditSpreadProfitLoss(tradeGroup) / tradeGroup.Credit) * 100;

      case 'Call Credit Spread':
        return (this.getCallCreditSpreadProfitLoss(tradeGroup) / tradeGroup.Credit) * 100;          
    }

    return 0.00
  } 

  //
  // Get trade progress for the bar
  //
  getTradeProgressBar(tradeGroup: TradeGroup) : number
  {    
    let p = 0.00;

    // Get progress based on type.
    switch(tradeGroup.Type)
    {
      case 'Put Credit Spread':
        p = (this.getPutCreditSpreadProfitLoss(tradeGroup) / tradeGroup.Credit) * 100;
      break;

      case 'Call Credit Spread':
        p = (this.getCallCreditSpreadProfitLoss(tradeGroup) / tradeGroup.Credit) * 100;
        break;

      case 'Option':
        p = this.getProgressOptionsTrade(tradeGroup);
      break;             
    }

    // keep it within a range
    if((p > 0) && (p <= 100))
    {
      return p;
    } else if(p > 100)
    {
      return 100;
    }    

    return 0.00
  } 

  //
  // Get the total P&L for put credit spreads
  //
  getPutCreditSpreadProfitLoss(tradeGroup: TradeGroup): number 
  {
    if(typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined")
    {
      return 0.00;
    }

    if(typeof this.quotes[tradeGroup.Positions[1].Symbol.ShortName] == "undefined")
    {
      return 0.00;
    }

    return tradeGroup.Credit - (((this.quotes[tradeGroup.Positions[1].Symbol.ShortName].ask - this.quotes[tradeGroup.Positions[0].Symbol.ShortName].bid) * 100) * Math.abs(tradeGroup.Positions[0].Qty));  
  }

  //
  // Get the total P&L for put credit spreads - Call Credit Spread
  //
  getCallCreditSpreadProfitLoss(tradeGroup: TradeGroup): number 
  {
    if (typeof this.quotes[tradeGroup.Positions[0].Symbol.ShortName] == "undefined") {
      return 0.00;
    }

    if (typeof this.quotes[tradeGroup.Positions[1].Symbol.ShortName] == "undefined") {
      return 0.00;
    }

    return tradeGroup.Credit - ((((this.quotes[tradeGroup.Positions[1].Symbol.ShortName].ask - this.quotes[tradeGroup.Positions[0].Symbol.ShortName].bid) * 100) * Math.abs(tradeGroup.Positions[0].Qty)) * -1);
  }

  //
  // Get trade profit and loss
  //
  getTradeProfitLoss(tradeGroup: TradeGroup) : number 
  {
    // Get progress based on type.
    switch(tradeGroup.Type)
    {
      case 'Put Credit Spread':
        return this.getPutCreditSpreadProfitLoss(tradeGroup);

      case 'Call Credit Spread':
        return this.getCallCreditSpreadProfitLoss(tradeGroup);        

      case 'Option':
        return this.getOptionProfitLoss(tradeGroup.Positions[0]);            
    }

    return 0.00;
  }

  //
  // Do we show a progress bar
  //
  showProgressbar() : boolean
  {
    return true;
  }  
}

/* End File */