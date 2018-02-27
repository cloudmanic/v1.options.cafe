//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit } from '@angular/core';
import { Balance } from '../../models/balance';
import { MarketStatus } from '../../models/market-status';
import { Broker } from '../../models/broker';
import { BrokerAccount } from '../../models/broker-account';
import { BrokerService } from '../../providers/http/broker.service';
import { StateService } from '../../providers/state/state.service';
import { WebsocketService } from '../../providers/http/websocket.service';

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

  private destory: Subject<boolean> = new Subject<boolean>();

  //
  // Construct.
  //
  constructor(private websocketService: WebsocketService, private brokerService: BrokerService, private stateService: StateService) { }

  //
  // Oninit...
  //
  ngOnInit() {
    this.brokerAccountList = [];
          
    // Subscribe to data updates from the broker - Market Status
    this.websocketService.marketStatusPush.subscribe(data => {
      this.marketStatus = data;      
    });

    // Subscribe to changes in the selected broker.
    this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
      this.getBalances();
    });

    // Subscribe to data updates from the broker - Balances
    this.websocketService.balancesPush.subscribe(data => {
      this.doBalanaces(data);
    });  

    // Get broker data
    this.brokerService.get().subscribe((data) => {
      this.doBrokers(data);
    });        
  }

  //
  // OnDestroy
  //
  ngOnDestroy()
  {
    this.destory.next();
    this.destory.complete();
  }   

  //
  // Get balances
  //
  getBalances()
  {
    // Get balance data
    this.brokerService.getBalances().subscribe((data) => {    
      this.doBalanaces(data);
    });
  }

  //
  // Do brokers
  //
  doBrokers(data: Broker[])
  {
    this.brokerList = data;

    let activeAccountId: string = this.stateService.GetStoredActiveAccountId();

    // Default to first one.
    if(! activeAccountId)
    {
      // Make sure we have at least one broker.
      if(! this.brokerList[0])
      {
        return false;
      }

      // Do we have a stored broker
      this.stateService.SetActiveBrokerAccount(this.brokerList[0].BrokerAccounts[0]);
      activeAccountId = this.stateService.GetStoredActiveAccountId();
    }

    // Loop through all the brokers and set our active broker. And make a list.
    for(var k = 0; k < this.brokerList.length; k++)
    {
      for(var i = 0; i < this.brokerList[k].BrokerAccounts.length; i++)
      {
        this.brokerAccountList.push(this.brokerList[k].BrokerAccounts[i]);

        // Set the selected account.
        if(String(this.brokerList[k].BrokerAccounts[i].Id) == activeAccountId)
        {
          this.selectedAccount = this.brokerList[k].BrokerAccounts[i];

          // Force refresh of balances
          this.stateService.SetActiveBrokerAccount(this.selectedAccount);           
        }
      }  
    }    
  }

  //
  // Do balances 
  //
  doBalanaces(data: Balance[])
  {
    // We have not gotten our brokers yet.
    if(! this.stateService.GetActiveBrokerAccount())
    {
      return false;
    }

    for(var i = 0; i < data.length; i++)
    {
      if(data[i].AccountNumber == this.stateService.GetActiveBrokerAccount().AccountNumber)
      {
        this.balance = data[i];
      }
    }
  }
  
  //
  // On account change.
  //
  onAccountChange() {
    this.stateService.SetActiveBrokerAccount(this.selectedAccount);
  }  
}

/* End File */