//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { AppService } from '../../providers/websocket/app.service';
import { Balance } from '../../models/balance';
import { UserProfile } from '../../models/user-profile';
import { MarketStatus } from '../../models/market-status';
import { BrokerAccounts } from '../../models/broker-accounts';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html'
})
export class SidebarComponent implements OnInit {

  balance: Balance;
  userProfile: UserProfile;
  marketStatus: MarketStatus;
  selectedAccount: BrokerAccounts;

  //
  // Construct.
  //
  constructor(private app: AppService, private changeDetect: ChangeDetectorRef) { }

  //
  // Oninit...
  //
  ngOnInit() {
    
    // Subscribe to data updates from the broker - User Profile
    this.app.userProfilePush.subscribe(data => {
      this.userProfile = data;
      this.changeDetect.detectChanges();
    });
    
    // Subscribe to data updates from the broker - Market Status
    this.app.userProfilePush.subscribe(data => {
      this.userProfile = data;
      
      // Do we have an account already? Always have to reset the selected one when we get new account data.
      if((! this.selectedAccount) && (this.userProfile.Accounts.length))
      {
        this.selectedAccount = this.userProfile.Accounts[0];
        this.app.setActiveAccountId(this.selectedAccount.AccountNumber);
      } else
      {
        for(var i = 0; i < this.userProfile.Accounts.length; i++)
        {
          if(this.userProfile.Accounts[i].AccountNumber == this.selectedAccount.AccountNumber)
          {
            this.selectedAccount = this.userProfile.Accounts[i];            
          }
        }
      }
      
      this.changeDetect.detectChanges();
    });    
        
  }
  
  //
  // On account change.
  //
  onAccountChange() {
    
    this.app.setActiveAccountId(this.selectedAccount.AccountNumber);
  
  }  

}

/* End File */