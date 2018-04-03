//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { TradeService, TradeEvent, TradeDetails } from '../../providers/http/trade.service';
import { Component, OnInit, ChangeDetectorRef } from '@angular/core';

@Component({
  selector: 'app-trading-sub-nav',
  templateUrl: './subnav.component.html'
})

export class SubnavComponent implements OnInit 
{
  //
  // Construct.
  //
  constructor(private tradeService: TradeService) { }

  //
  // OnInit
  //
  ngOnInit() 
  {
  }

  //
  // On Trade Click.
  //
  onTradeClick()
  {
    this.tradeService.PushEvent(new TradeEvent().createNew("toggle-trade-builder", new TradeDetails()));
  }

}

/* End File */