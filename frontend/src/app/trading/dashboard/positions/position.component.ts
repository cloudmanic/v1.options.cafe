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
  getOptionWidgetProfitLoss(tradeGroups: TradeGroup) : number 
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

}

/* End File */