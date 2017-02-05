import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { BrokerService } from '../services/broker.service';
import { Order } from '../contracts/order';

@Component({
  selector: 'oc-orders',
  templateUrl: './orders.component.html'
})
export class OrdersComponent implements OnInit {
  
  orders: Order[]
  activeAccount = ""

  //
  // Constructor....
  //
  constructor(private broker: BrokerService, private changeDetect: ChangeDetectorRef) { }

  //
  // OnInit....
  //
  ngOnInit() {
    
    // Set the active account.
    this.activeAccount = this.broker.activeAccount;
    
    // Subscribe to data updates from the broker - Orders
    this.broker.ordersPushData.subscribe(data => {
      
      var rt = []
      
      // Filter - We only one the accounts that are active.
      for(var i = 0; i < data.length; i++)
      {        
        if(data[i].AccountId == this.activeAccount)
        {
          rt.push(data[i]);
        }
      }
      
      // Set order data
      this.orders = rt;      
      this.changeDetect.detectChanges();
    });    
    
    // Subscribe to when the active account changes
    this.broker.activeAccountPushData.subscribe(data => {
      this.activeAccount = data;
      this.orders = [];
      this.changeDetect.detectChanges();
    }); 
    
  }

}

/* End File */