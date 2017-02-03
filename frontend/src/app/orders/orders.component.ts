import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { BrokerService } from '../services/broker.service';
import { Order } from '../contracts/order';

@Component({
  selector: 'oc-orders',
  templateUrl: './orders.component.html'
})
export class OrdersComponent implements OnInit {
  
  orders: Order[]

  //
  // Constructor....
  //
  constructor(private broker: BrokerService, private changeDetect: ChangeDetectorRef) { }

  //
  // OnInit....
  //
  ngOnInit() {
    
    // Subscribe to data updates from the broker - Orders
    this.broker.ordersPushData.subscribe(data => {
      this.orders = data;      
      this.changeDetect.detectChanges();
    });    
    
  }

}

/* End File */