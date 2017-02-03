import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { BrokerService } from '../services/broker.service';
import { UserProfile } from '../contracts/user-profile';
import { MarketStatus } from '../contracts/market-status';

@Component({
  selector: 'oc-sidebar',
  templateUrl: './sidebar.component.html'
})
export class SidebarComponent implements OnInit {

  userProfile: UserProfile;
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
    this.broker.userProfilePushData.subscribe(data => {
      this.userProfile = data;
      this.changeDetect.detectChanges();
    });
    
    // Subscribe to data updates from the broker - Market Status
    this.broker.marketStatusPushData.subscribe(data => {
      this.marketStatus = data;      
      this.changeDetect.detectChanges();
    });
     
  }
}

/* End File */