//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { WebsocketService } from '../../providers/http/websocket.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './home.component.html'
})
export class DashboardComponent implements OnInit {

  ws_reconnecting = false;

  //
  // Construct...
  //
  constructor(private websocketService: WebsocketService, private changeDetect: ChangeDetectorRef) { }

  //
  // On Init...
  //
  ngOnInit() {
    
    // Subscribe to when we are reconnecting to a websocket - Core
    this.websocketService.wsReconnecting.subscribe(data => {
      this.ws_reconnecting = data;
      this.changeDetect.detectChanges();
    });

  }

}

/* End File */