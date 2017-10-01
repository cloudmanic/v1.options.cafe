//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { AppService } from '../../providers/websocket/app.service';
import { QuoteService } from '../../providers/websocket/quote.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './home.component.html'
})
export class DashboardComponent implements OnInit {

  ws_reconnecting = false;

  //
  // Construct...
  //
  constructor(private appService: AppService, private quoteService: QuoteService, private changeDetect: ChangeDetectorRef) { }

  //
  // On Init...
  //
  ngOnInit() {
    
    // Subscribe to when we are reconnecting to a websocket - Core
    this.appService.wsReconnecting.subscribe(data => {
      this.ws_reconnecting = data;
      this.changeDetect.detectChanges();
    });     

    // Subscribe to when we are reconnecting to a websocket - Quotes
    this.quoteService.wsReconnecting.subscribe(data => {
      this.ws_reconnecting = data;
      this.changeDetect.detectChanges();
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