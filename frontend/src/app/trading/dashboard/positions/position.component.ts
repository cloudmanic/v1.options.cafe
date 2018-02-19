//
// Date: 2/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { TradeGroup } from '../../../models/trade-group';
import { Position } from '../../../models/position';
import { Component, OnInit, Input, Output, OnChanges, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-trading-position',
  templateUrl: './position.component.html'
})

export class PositionComponent implements OnInit 
{
  @Input() quotes = {};
  @Input() title: string;   
  @Input() tradeGroups: TradeGroup[]; 

  //
  // Constructor....
  //
  constructor() { }

  //
  // OnInit....
  //
  ngOnInit() 
  { 
    console.log(this.tradeGroups);

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
  // Get trade progress
  //
  getTradeProgress(tradeGroup: TradeGroup) : number
  {    
    switch(tradeGroup.Type)
    {
      case 'Put Credit Spread':
        return (this.getPutCreditSpreadProfitLoss(tradeGroup) / tradeGroup.Credit) * 100;
      break;     
    }

    return 0.00
  } 

  //
  // Get trade progress for the bar
  //
  getTradeProgressBar(tradeGroup: TradeGroup) : number
  {    
    switch(tradeGroup.Type)
    {
      case 'Put Credit Spread':
        let p = (this.getPutCreditSpreadProfitLoss(tradeGroup) / tradeGroup.Credit) * 100;
        if(p > 0)
        {
          return p;
        }
      break;     
    }

    return 0.00
  } 

  //
  // Get the total P&L for put credit spreads
  //
  getPutCreditSpreadProfitLoss(tradeGroup: TradeGroup) : number 
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
}

    // if(! $scope.quotes[spread.Positions[1].SymbolsShort])
    // {
    //   return 0;
    // }
    
    // if(type == 'put')
    // {
      
    //   //console.log((spread.TradeGroupsOpen * -1) - ((($scope.quotes[spread.Positions[1].SymbolsShort].ask - $scope.quotes[spread.Positions[0].SymbolsShort].bid) * 100) * spread.Positions[0].PositionsQty) );
      
    //   return (spread.TradeGroupsOpen * -1) - ((($scope.quotes[spread.Positions[1].SymbolsShort].ask - $scope.quotes[spread.Positions[0].SymbolsShort].bid) * 100) * spread.Positions[0].PositionsQty)       
    // } else
    // {
    //   return (spread.TradeGroupsOpen * -1) - ((($scope.quotes[spread.Positions[0].SymbolsShort].ask - $scope.quotes[spread.Positions[1].SymbolsShort].bid) * 100) * spread.Positions[1].PositionsQty)      
    // }

/* End File */