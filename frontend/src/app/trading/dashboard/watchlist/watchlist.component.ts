//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AppService } from '../../../providers/websocket/app.service';
import { QuoteService } from '../../../providers/websocket/quote.service';
import { Watchlist } from '../../../models/watchlist';

@Component({
  selector: 'app-watchlist',
  templateUrl: './watchlist.component.html'
})

export class WatchlistComponent implements OnInit {

  public quotes = {}
  public watchlist: Watchlist;

  //
  // Construct...
  //
  constructor(private appService: AppService, private quoteService: QuoteService) { }

  //
  // On Init...
  //
  ngOnInit() {
    
    this.watchlist = this.appService.watchlist;
    
    // Subscribe to data updates from the backend - Watchlist
    this.appService.watchlistPush.subscribe(data => {      
      this.watchlist = data;
    });    
    
    // Subscribe to data updates from the quotes - Market Quotes
    this.quoteService.marketQuotePushData.subscribe(data => {
      this.quotes[data.symbol] = data;
    });   
 
  }

  //
  // On watchlist settings click.
  //
  onWatchlistSettingsClick() {
    this.appService.RequestWatchlistData();
  } 

}

/* End File */