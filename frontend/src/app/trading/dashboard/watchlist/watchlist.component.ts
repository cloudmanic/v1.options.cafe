//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ElementRef } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { AppService } from '../../../providers/websocket/app.service';
import { QuoteService } from '../../../providers/websocket/quote.service';
import { Symbol } from '../../../models/symbol';
import { Watchlist } from '../../../models/watchlist';
import { environment } from '../../../../environments/environment';

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

  //
  // Construct...
  //
  constructor(private http: HttpClient, private _eref: ElementRef, private appService: AppService, private quoteService: QuoteService) { }

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
  // onSearchTypeAheadClick() 
  //
  onSearchTypeAheadClick(symbol: Symbol) {
   console.log(symbol)
  }

  //
  // On watchlist settings click.
  //
  onWatchlistSettingsClick() {

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
  // Do edit watchlist
  //
  onEditWatchList() {
    this.watchlistEditState = true;
    this.watchlistSettingsActive = false;
  }

  //
  // Do static watchlist
  //
  onEditWatchListDone() {
    this.watchlistEditState = false;
    this.watchlistSettingsActive = false;
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