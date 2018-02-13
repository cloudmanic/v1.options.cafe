//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AppService } from '../../../providers/websocket/app.service';
import { QuoteService } from '../../../providers/websocket/quote.service';
import { Order } from '../../../models/order';
import { BrokerStateService } from '../../../providers/state/broker.state.service';

@Component({
  selector: 'app-trading-orders',
  templateUrl: './orders.component.html'
})

export class OrdersComponent implements OnInit {
  
  quotes = {}
  orders: Order[]  

  //
  // Constructor....
  //
  constructor(private appService: AppService, private quoteService: QuoteService, private brokerState: BrokerStateService) { }

  //
  // OnInit....
  //
  ngOnInit() {
    // Get Data from cache
    this.setOrders(this.appService.orders);
    this.quotes = this.quoteService.quotes;
        
    // Subscribe to data updates from the broker - Orders
    this.appService.ordersPush.subscribe(data => {
      this.setOrders(data);
    });    
    
    // Subscribe to data updates from the quotes - Market Quotes
    this.quoteService.marketQuotePushData.subscribe(data => {
      this.quotes[data.symbol] = data;
    });     
  }

  //
  // Set the orders.
  //
  private setOrders(orders: Order[]) {
    var rt = []
    
    // This data has not come in yet.
    if(! this.brokerState.GetActiveBrokerAccount())
    {
      return;
    }      
    
    // Filter - We only one the accounts that are active.
    for(var i = 0; i < orders.length; i++)
    {                
      if(orders[i].AccountId == this.brokerState.GetActiveBrokerAccount().AccountNumber)
      {
        rt.push(orders[i]);
      }
    }
    
    // Set order data
    this.orders = rt;
  }
}

/* End File */