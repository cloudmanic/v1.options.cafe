//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ElementRef } from '@angular/core';
import { AppService } from '../../../providers/websocket/app.service';
import { QuoteService } from '../../../providers/websocket/quote.service';
import { Watchlist } from '../../../models/watchlist';

@Component({
  selector: 'app-watchlist',
  templateUrl: './watchlist.component.html',
  host: { '(document:click)': 'onDocClick($event)' }
})

export class WatchlistComponent implements OnInit {

  public quotes = {}
  public watchlist: Watchlist;
  public watchlistEditState = false;
  public watchlistSettingsActive = false;

  public typeAheadList = [
    { symbol: 'spy', description: 'SPDR S&P 500 ETF Trust' },
    { symbol: 'sbux', description: 'Starbucks Corp' },
    { symbol: 'bac', description: 'Bank of America' }       
  ];

  //
  // Construct...
  //
  constructor(private _eref: ElementRef, private appService: AppService, private quoteService: QuoteService) { }

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

    this.appService.RequestSymbolSearch('sbux');

    if(this.watchlistSettingsActive)
    {
      this.watchlistSettingsActive = false;
    } else
    {
      this.watchlistSettingsActive = true;      
    }

    //this.appService.RequestWatchlistData();
  } 

  //
  // Click anywhere on the screen.
  //
  onDocClick(event) {
   
     // Remove active buttons
    if(! this._eref.nativeElement.contains(event.target))
    {
      this.watchlistSettingsActive = false;
    }

  }  

}

/* End File */