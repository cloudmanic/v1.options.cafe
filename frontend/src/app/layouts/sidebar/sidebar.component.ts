//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { AppState } from '../../providers/app.state.service';
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
  constructor(private app: AppService, private appState: AppState, private changeDetect: ChangeDetectorRef) { }

  //
  // Oninit...
  //
  ngOnInit() {
    
/*
    // Pull in from the cache the default values.
    this.userProfile = this.appState.getUserProfile();    
    this.selectedAccount = this.appState.getActiveAccount();
    this.changeDetect.detectChanges();
    
    console.log(this.userProfile);
*/
        
    // Subscribe to data updates from the broker - Market Status
    this.app.marketStatusPush.subscribe(data => {
      this.marketStatus = data;      
    });    
            
    // Subscribe to data updates from the broker - Market Status
    this.app.userProfilePush.subscribe(data => {
      
      this.userProfile = data;
      this.appState.setUserProfile(data);
      
      // Do we have an account already? Always have to reset the selected one when we get new account data.
      if((! this.selectedAccount) && (this.userProfile.Accounts.length))
      {
        this.selectedAccount = this.userProfile.Accounts[0];
        this.app.setActiveAccount(this.selectedAccount);
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
      
    }); 
    
    // Subscribe to data updates from the broker - Balances
    this.app.balancesPush.subscribe(data => {

      for(var i = 0; i < data.length; i++)
      {
        if(data[i].AccountNumber == this.app.activeAccount)
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