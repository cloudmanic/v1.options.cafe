//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AppService } from '../../providers/websocket/app.service';
import { Balance } from '../../models/balance';
import { UserProfile } from '../../models/user-profile';
import { MarketStatus } from '../../models/market-status';
import { BrokerAccount } from '../../models/broker-account';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html'
})
export class SidebarComponent implements OnInit {

  balance: Balance;
  userProfile: UserProfile;
  marketStatus: MarketStatus;
  selectedAccount: BrokerAccount;

  //
  // Construct.
  //
  constructor(private app: AppService) { }

  //
  // Oninit...
  //
  ngOnInit() {
          
    // Subscribe to data updates from the broker - Market Status
    this.app.marketStatusPush.subscribe(data => {
      this.marketStatus = data;      
    });    
            
    // Subscribe to data updates from the broker - Market Status
    this.app.userProfilePush.subscribe(data => {
      
      this.userProfile = data;

      if(! this.userProfile.Accounts.length)
      {
        return;
      }

      if(! this.app.getActiveAccount())
      {
        this.selectedAccount = data.Accounts[0];
        return;
      }
      
      for(var i = 0; i < this.userProfile.Accounts.length; i++)
      {
        if(this.userProfile.Accounts[i].AccountNumber == this.app.getActiveAccount().AccountNumber)
        {
          this.selectedAccount = this.userProfile.Accounts[i];           
        }
      }
      
    }); 
    
    // Subscribe to data updates from the broker - Balances
    this.app.balancesPush.subscribe(data => {

      for(var i = 0; i < data.length; i++)
      {
        if(data[i].AccountNumber == this.app.getActiveAccount().AccountNumber)
        {
          this.balance = data[i];
        }
      }

    });       
        
  }
  
  //
  // On account change.
  //
  onAccountChange() {
    this.app.setActiveAccount(this.selectedAccount);
  }  

}

/* End File */