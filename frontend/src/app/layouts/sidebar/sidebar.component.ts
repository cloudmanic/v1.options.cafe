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
import { Broker } from '../../models/broker';
import { BrokerAccount } from '../../models/broker-account';
import { BrokerService } from '../../providers/http/broker.service';
import { BrokerStateService } from '../../providers/state/broker.state.service';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html'
})
export class SidebarComponent implements OnInit {

  brokerList: Broker[];
  balance: Balance;
  userProfile: UserProfile;
  marketStatus: MarketStatus;
  selectedAccount: BrokerAccount;
  brokerAccountList: BrokerAccount[];

  //
  // Construct.
  //
  constructor(private app: AppService, private brokerService: BrokerService, private brokerState: BrokerStateService) { }

  //
  // Oninit...
  //
  ngOnInit() {

    this.brokerAccountList = [];

    this.getBrokers()
          
    // Subscribe to data updates from the broker - Market Status
    this.app.marketStatusPush.subscribe(data => {
      this.marketStatus = data;      
    });    
            
    // // Subscribe to data updates from the broker - Market Status
    // this.app.userProfilePush.subscribe(data => {
      
    //   this.userProfile = data;

    //   if(! this.userProfile.Accounts.length)
    //   {
    //     return;
    //   }

    //   if(! this.app.getActiveAccount())
    //   {
    //     this.selectedAccount = data.Accounts[0];
    //     return;
    //   }
      
    //   for(var i = 0; i < this.userProfile.Accounts.length; i++)
    //   {
    //     if(this.userProfile.Accounts[i].AccountNumber == this.app.getActiveAccount().AccountNumber)
    //     {
    //       this.selectedAccount = this.userProfile.Accounts[i];           
    //     }
    //   }
      
    // }); 
    
    // Subscribe to data updates from the broker - Balances
    this.app.balancesPush.subscribe(data => {

      // We have not gotten our brokers yet.
      if(! this.brokerState.GetActiveBrokerAccount())
      {
        return false;
      }

      for(var i = 0; i < data.length; i++)
      {
        if(data[i].AccountNumber == this.brokerState.GetActiveBrokerAccount().AccountNumber)
        {
          this.balance = data[i];
        }
      }

    });       
        
  }

  //
  // Get brokers
  //
  getBrokers() {

    // Get broker data
    this.brokerService.get().subscribe((data) => {
      this.brokerList = data;

      // Default to first one.
      if(! this.brokerState.GetActiveBrokerAccount())
      {
        // Make sure we have at least one broker.
        if(! this.brokerList[0])
        {
          return false;
        }

        this.brokerState.SetActiveBrokerAccount(this.brokerList[0].BrokerAccounts[0]);
      }

      // Loop through all the brokers and set our active broker. And make a list.
      for(var k = 0; k < this.brokerList.length; k++)
      {
        for(var i = 0; i < this.brokerList[k].BrokerAccounts.length; i++)
        {
          this.brokerAccountList.push(this.brokerList[k].BrokerAccounts[i]);

          // Set the selected account.
          if(this.brokerList[k].BrokerAccounts[i].Id == this.brokerState.GetActiveBrokerAccount().Id)
          {
            this.selectedAccount = this.brokerList[k].BrokerAccounts[i];           
          }
        }  
      }
    });    
  }  
  
  //
  // On account change.
  //
  onAccountChange() {
    this.brokerState.SetActiveBrokerAccount(this.selectedAccount);
  }  
}

/* End File */