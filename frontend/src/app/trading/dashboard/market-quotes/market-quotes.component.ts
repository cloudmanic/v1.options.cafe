//
// Date: 2/21/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AppService } from '../../../providers/websocket/app.service';
import { QuoteService } from '../../../providers/websocket/quote.service';
import { StateService } from '../../../providers/state/state.service';

@Component({
  selector: 'app-trading-market-quotes',
  templateUrl: './market-quotes.component.html'
})

export class MarketQuotesComponent implements OnInit {
  
  private quotes = {}

  //
  // Constructor....
  //
  constructor(private appService: AppService, private quoteService: QuoteService, private stateService: StateService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Get Data from cache
    this.quotes = this.quoteService.quotes;
            
    // Subscribe to data updates from the quotes - Market Quotes
    this.quoteService.marketQuotePushData.subscribe(data => {
      this.quotes[data.symbol] = data;
    });     
  }  

}

/* End File */