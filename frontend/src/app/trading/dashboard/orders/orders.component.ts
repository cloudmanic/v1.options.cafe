//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AppService } from '../../../providers/websocket/app.service';
import { QuoteService } from '../../../providers/websocket/quote.service';
import { Order } from '../../../models/order';
import { BrokerAccount } from '../../../models/broker-account';

@Component({
  selector: 'app-trading-orders',
  templateUrl: './orders.component.html'
})

export class OrdersComponent implements OnInit {
  
  quotes = {}
  orders: Order[]  
  activeAccount: BrokerAccount

  //
  // Constructor....
  //
  constructor(private app: AppService, private quoteService: QuoteService) { }

  //
  // OnInit....
  //
  ngOnInit() {
        
    // Subscribe to data updates from the broker - Orders
    this.app.ordersPush.subscribe(data => {
      
      var rt = []
      
      // Set the active account.
      this.activeAccount = this.app.getActiveAccount();
      
      // This data has not come in yet.
      if(! this.activeAccount)
      {
        return;
      }      
      
      // Filter - We only one the accounts that are active.
      for(var i = 0; i < data.length; i++)
      {                
        if(data[i].AccountId == this.activeAccount.AccountNumber)
        {
          rt.push(data[i]);
        }
      }
      
      // Set order data
      this.orders = rt;
      
    });    
    
    // Subscribe to data updates from the quotes - Market Quotes
    this.quoteService.marketQuotePushData.subscribe(data => {
    
      this.quotes[data.symbol] = data;
      //this.changeDetect.detectChanges();
      
    });     
    
/*
    // Subscribe to when the active account changes
    this.broker.activeAccountPushData.subscribe(data => {
      this.activeAccount = data;
      this.changeDetect.detectChanges();
    });
*/    
    
  }

}

/* End File */