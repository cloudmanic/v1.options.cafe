//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { AppService } from '../../providers/websocket/app.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './home.component.html'
})
export class DashboardComponent implements OnInit {

  ws_reconnecting = false;

  //
  // Construct...
  //
  constructor(private app: AppService, private changeDetect: ChangeDetectorRef) { }

  //
  // On Init...
  //
  ngOnInit() {
    
    // Subscribe to when we are reconnecting to a websocket - Core
    this.app.wsReconnecting.subscribe(data => {
      this.ws_reconnecting = data;
      this.changeDetect.detectChanges();
    });     

    // Subscribe to when we are reconnecting to a websocket - Quotes
    this.app.wsReconnecting.subscribe(data => {
      this.ws_reconnecting = data;
      this.changeDetect.detectChanges();
    });    
    
  }

}
