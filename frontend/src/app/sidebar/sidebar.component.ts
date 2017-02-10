import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { Balance } from '../contracts/balance';
import { BrokerService } from '../services/broker.service';
import { UserProfile } from '../contracts/user-profile';
import { MarketStatus } from '../contracts/market-status';

@Component({
  selector: 'oc-sidebar',
  templateUrl: './sidebar.component.html'
})
export class SidebarComponent implements OnInit {
  
  balance: Balance

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
  
    // Subscribe to data updates from the broker - User Profile
    this.broker.userProfilePushData.subscribe(data => {
      this.userProfile = data;
      this.changeDetect.detectChanges();
    });
    
    // Subscribe to data updates from the broker - Market Status
    this.broker.marketStatusPushData.subscribe(data => {
      this.marketStatus = data;      
      this.changeDetect.detectChanges();
    });
    
    // Subscribe to data updates from the broker - Balances
    this.broker.balancesPushData.subscribe(data => {

      for(var i = 0; i < data.length; i++)
      {
        if(data[i].AccountNumber == this.broker.activeAccount)
        {
          this.balance = data[i];
        }
      }
    
      this.changeDetect.detectChanges();
    });    
     
  }
}

/* End File */