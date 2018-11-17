//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { AnalyzeService } from '../../providers/http/analyze.service'
import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { WebsocketService } from '../../providers/http/websocket.service';
import { NotificationsService } from '../../providers/http/notifications.service';
import { AnalyzeTrade } from '../../providers/http/analyze.service'
import { AnalyzeLeg } from '../../models/analyze-result';

@Component({
  selector: 'app-dashboard',
  templateUrl: './home.component.html'
})

export class DashboardComponent implements OnInit 
{
  openPrice: number = 88.99;
  showNotice: boolean = false;
  noticeId: number = 0;
  noticeTitle: string = '';
  noticeBody: string = '';
  ws_reconnecting: boolean = false;

  //
  // Construct...
  //
  constructor(private notificationsService: NotificationsService, private websocketService: WebsocketService, private changeDetect: ChangeDetectorRef, private analyzeService: AnalyzeService) 
  { 
    // Load data for page.
    this.getNotifications();
  }

  //
  // On Init...
  //
  ngOnInit() 
  {  
    // Subscribe to when we are reconnecting to a websocket - Core
    this.websocketService.wsReconnecting.subscribe(data => {
      this.ws_reconnecting = data;
      this.changeDetect.detectChanges();
    });

  }

  blah()
  {
    let trade = new AnalyzeTrade();
    trade.OpenCost = 100.00;

    let leg1 = new AnalyzeLeg();
    leg1.Qty = 1;
    leg1.SymbolStr = "SPY181221C00250000";

    let leg2 = new AnalyzeLeg();
    leg2.Qty = -2;
    leg2.SymbolStr = "SPY181221C00260000";

    let leg3 = new AnalyzeLeg();
    leg3.Qty = 1;
    leg3.SymbolStr = "SPY181221C00270000";    
    
    trade.Legs = [ leg1, leg2, leg3 ];


    this.analyzeService.dialog.emit(trade);
  }

  //
  // See if we have any notifications to display.
  //
  getNotifications()
  {
    this.notificationsService.get("in-app", "dashboard-notice", "pending").subscribe(data => {

      // We only take one notice at a time.
      if(data.length <= 0)
      {
        this.noticeId = 0;
        this.showNotice = false;
        this.noticeTitle = "";
        this.noticeBody = "";        
        return;
      }

      // Grab the first notice and display.
      this.noticeId = data[0].Id;
      this.noticeTitle = data[0].Title;
      this.noticeBody = data[0].LongMessage;
      this.showNotice = true;

    });
  }

  //
  // Close Notice 
  //
  closeNotice()
  {
    // Send API call to mark notice as seen. 
    this.notificationsService.markSeen(this.noticeId).subscribe(data => {

      // See if there is a "next" notice.
      this.getNotifications();

    });
  }

}

/* End File */