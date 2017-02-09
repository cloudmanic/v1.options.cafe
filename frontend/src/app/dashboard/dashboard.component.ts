import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { QuoteService } from '../services/quote.service';
import { BrokerService } from '../services/broker.service';

@Component({
  selector: 'oc-dashboard',
  templateUrl: './dashboard.component.html'
})
export class DashboardComponent implements OnInit {

  ws_reconnecting = false;

  constructor(private quotesService: QuoteService, private broker: BrokerService, private changeDetect: ChangeDetectorRef) { }

  //
  // OnInit....
  //
  ngOnInit() {
    
    // Subscribe to when we are reconnecting to a websocket - Core
    this.broker.wsReconnecting.subscribe(data => {
      this.ws_reconnecting = data;
      this.changeDetect.detectChanges();
    });     

    // Subscribe to when we are reconnecting to a websocket - Quotes
    this.quotesService.wsReconnecting.subscribe(data => {
      this.ws_reconnecting = data;
      this.changeDetect.detectChanges();
    }); 
    
  }

}
