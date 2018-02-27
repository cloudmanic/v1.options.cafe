//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit } from '@angular/core';
import { Order } from '../../../models/order';
import { StateService } from '../../../providers/state/state.service';
import { WebsocketService } from '../../../providers/http/websocket.service';

@Component({
  selector: 'app-trading-orders',
  templateUrl: './orders.component.html'
})

export class OrdersComponent implements OnInit {
  
  quotes = {}
  orders: Order[]  

  private destory: Subject<boolean> = new Subject<boolean>();  

  //
  // Constructor....
  //
  constructor(private websocketService: WebsocketService, private stateService: StateService) { }

  //
  // OnInit....
  //
  ngOnInit() {
    // Get Data from cache
    this.quotes = this.stateService.GetQuotes();
        
    // Subscribe to data updates from the broker - Orders
    this.websocketService.ordersPush.subscribe(data => {
      this.setOrders(data);
    });    
    
    // Subscribe to data updates from the quotes - Market Quotes
    this.websocketService.quotePushData.subscribe(data => {
      this.quotes[data.symbol] = data;
    });     
  
    // Subscribe to changes in the selected broker.
    this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
      this.orders = null;
    });
  }

  //
  // OnDestroy
  //
  ngOnDestroy()
  {
    this.destory.next();
    this.destory.complete();
  }  

  //
  // Set the orders.
  //
  private setOrders(orders: Order[]) {
    var rt = []
    
    // This data has not come in yet.
    if(! this.stateService.GetActiveBrokerAccount())
    {
      return;
    }      
    
    // Filter - We only one the accounts that are active.
    for(var i = 0; i < orders.length; i++)
    {                
      if(orders[i].AccountId == this.stateService.GetActiveBrokerAccount().AccountNumber)
      {
        rt.push(orders[i]);
      }
    }
    
    // Set order data
    this.orders = rt;
  }
}

/* End File */