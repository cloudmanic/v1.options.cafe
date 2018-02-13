//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AppService } from '../../providers/websocket/app.service';
import { Balance } from '../../models/balance';
import { MarketStatus } from '../../models/market-status';
import { Broker } from '../../models/broker';
import { BrokerAccount } from '../../models/broker-account';
import { BrokerService } from '../../providers/http/broker.service';
import { StateService } from '../../providers/state/state.service';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html'
})
export class SidebarComponent implements OnInit {

  brokerList: Broker[];
  balance: Balance;
  marketStatus: MarketStatus;
  selectedAccount: BrokerAccount;
  brokerAccountList: BrokerAccount[];

  //
  // Construct.
  //
  constructor(private app: AppService, private brokerService: BrokerService, private siteService: StateService) { }

  //
  // Oninit...
  //
  ngOnInit() {

    this.brokerAccountList = [];

    // API call to get brokers
    this.getBrokers()
          
    // Subscribe to data updates from the broker - Market Status
    this.app.marketStatusPush.subscribe(data => {
      this.marketStatus = data;      
    });
    
    // Subscribe to data updates from the broker - Balances
    this.app.balancesPush.subscribe(data => {

      // We have not gotten our brokers yet.
      if(! this.siteService.GetActiveBrokerAccount())
      {
        return false;
      }

      for(var i = 0; i < data.length; i++)
      {
        if(data[i].AccountNumber == this.siteService.GetActiveBrokerAccount().AccountNumber)
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

      let activeAccountId = this.siteService.GetStoredActiveAccountId();

      // Default to first one.
      if(! activeAccountId)
      {
        // Make sure we have at least one broker.
        if(! this.brokerList[0])
        {
          return false;
        }

        // Do we have a stored broker
        this.siteState.SetActiveBrokerAccount(this.brokerList[0].BrokerAccounts[0]);
        activeAccountId = this.siteService.GetStoredActiveAccountId();
      }

      // Loop through all the brokers and set our active broker. And make a list.
      for(var k = 0; k < this.brokerList.length; k++)
      {
        for(var i = 0; i < this.brokerList[k].BrokerAccounts.length; i++)
        {
          this.brokerAccountList.push(this.brokerList[k].BrokerAccounts[i]);

          // Set the selected account.
          if(this.brokerList[k].BrokerAccounts[i].Id == activeAccountId)
          {
            this.selectedAccount = this.brokerList[k].BrokerAccounts[i];

            // Force refresh of balances
            this.siteService.SetActiveBrokerAccount(this.selectedAccount);           
          }
        }  
      }
    });    
  }  
  
  //
  // On account change.
  //
  onAccountChange() {
    this.siteService.SetActiveBrokerAccount(this.selectedAccount);
  }  
}

/* End File */