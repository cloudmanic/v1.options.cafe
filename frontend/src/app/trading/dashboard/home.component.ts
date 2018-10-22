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

export class DashboardComponent implements OnInit 
{
  showNotice: boolean = true;
  noticeTitle: string = "Welcome to Options Cafe Beta";
  noticeBody: string = 'Lorem ipsum dolor sit amet, adipisicing elit. Aperiam ab id quos eos sapiente nostrum voluptatem impedit vitae repellat voluptate quam eius temporibus necessitatibus, ea eveniet molestias deserunt, suscipit magni, <a href="#">incidunt excepturi rem</a> voluptates soluta, officiis animi porro! Facilis, enim.';
  ws_reconnecting: boolean = false;

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