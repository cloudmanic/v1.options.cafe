//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Rx';
import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit } from '@angular/core';
import { Order } from '../../../models/order';
import { StateService } from '../../../providers/state/state.service';
import { TradeService } from '../../../providers/http/trade.service';
import { WebsocketService } from '../../../providers/http/websocket.service';
import { BrokerService } from '../../../providers/http/broker.service';
import { ChangeDetected } from '../../../models/change-detected';
import { DropdownAction } from '../../../shared/dropdown-select/dropdown-select.component';

@Component({
  selector: 'app-trading-orders',
  templateUrl: './orders.component.html'
})

export class OrdersComponent implements OnInit {
  
  quotes = {}
  orders: Order[]
  actions: DropdownAction[] = null;   

  private destory: Subject<boolean> = new Subject<boolean>();  

  //
  // Constructor....
  //
  constructor(private websocketService: WebsocketService, private stateService: StateService, private brokerService: BrokerService, private tradeService: TradeService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Setup Dropdown actions
    this.setupDropdownActions();

    // Set orders
    this.orders = this.stateService.GetActiveOrders();

    // Get Data from cache
    this.quotes = this.stateService.GetQuotes();
            
    // Subscribe to data updates from the quotes - Market Quotes
    this.websocketService.quotePushData.takeUntil(this.destory).subscribe(data => {
      this.quotes[data.symbol] = data;
    });     
  
    // Subscribe to data updates from the broker - Orders
    this.websocketService.ordersPush.takeUntil(this.destory).subscribe(data => {
      this.doOrders(data);
    }); 

    // Subscribe to changes in the selected broker.
    this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
      this.getOrders();
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
  // Setup Drop down actions.
  //
  setupDropdownActions() 
  {
    let das = []

    // First action
    let da1 = new DropdownAction();
    da1.title = "Cancel Trade";

    // Cancel order action
    da1.click = (row: Order) => {

      this.tradeService.cancelOrder(this.stateService.GetActiveBrokerAccount().Id, row.Id).subscribe((res) => {

        // Show success notice
        this.stateService.SiteSuccess.emit("Order Canceled: Your order number #" + row.Id);

      });

    };

    das.push(da1);

    this.actions = das;
  }  

  //
  // Get Orders - Used for broker change
  //
  getOrders()
  {
    // Get balance data
    this.brokerService.getOrders(this.stateService.GetActiveBrokerAccount().BrokerId).subscribe((data) => {
      this.doOrders(data);
    });
  }  

  //
  // Set the orders.
  //
  doOrders(orders: Order[]) 
  {
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
    
    // See if anything changed if not no need to update the UI
    if(JSON.stringify(rt) == JSON.stringify(this.orders))
    {
      return;
    }

    // Set order data
    this.orders = rt;

    // Add orders to state manager
    this.stateService.SetActiveOrders(this.orders);
  }
}

/* End File */