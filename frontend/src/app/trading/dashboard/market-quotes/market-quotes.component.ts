//
// Date: 2/21/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { WebsocketService } from '../../../providers/http/websocket.service';
import { StateService } from '../../../providers/state/state.service';

@Component({
  selector: '[app-trading-market-quotes]',
  templateUrl: './market-quotes.component.html'
})

export class MarketQuotesComponent implements OnInit {
  
  public quotes = {}

  //
  // Constructor....
  //
  constructor(private websocketService: WebsocketService, private stateService: StateService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Get Data from cache
    this.quotes = this.stateService.GetQuotes();
            
    // Subscribe to data updates from the quotes - Market Quotes
    this.websocketService.quotePushData.subscribe(data => {
      this.quotes[data.symbol] = data;
    });     
  }  

}

/* End File */