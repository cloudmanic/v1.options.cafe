import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { BrokerService } from '../services/broker.service';
import { MarketStatus } from '../contracts/market-status';

@Component({
  selector: 'oc-footer',
  templateUrl: './footer.component.html'
})
export class FooterComponent implements OnInit {
  marketStatus: MarketStatus;

  //
  // Constructor....
  //
  constructor(private broker: BrokerService, private changeDetect: ChangeDetectorRef) { }

  //
  // OnInit....
  //
  ngOnInit() {

    // Subscribe to data updates from the broker - Market Status
    this.broker.marketStatusPushData.subscribe(data => {
      this.marketStatus = data;
      this.changeDetect.detectChanges();
    });

  }

}
