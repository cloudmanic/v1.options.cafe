//
// Date: 2/13/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Injectable } from '@angular/core';
import { Broker } from '../../models/broker';
import { BrokerAccount } from '../../models/broker-account';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { environment } from '../../../environments/environment';
import { AppService } from '../../providers/websocket/app.service';

@Injectable()
export class BrokerStateService  
{ 
  private activeBrokerAccount: BrokerAccount

  //
  // Construct.
  //
  constructor(private appService: AppService) { 
    this.activeBrokerAccount = null;
  }

  //
  // Get stored active account id
  //
  GetStoredActiveAccountId() : number {
    return localStorage.getItem('active_account')
  }

  //
  // Set Active Broker Account
  //
  SetActiveBrokerAccount(brokerAccount: BrokerAccount) {
    this.activeBrokerAccount = brokerAccount
    localStorage.setItem('active_account', brokerAccount.Id);
    this.appService.RequestAllData();
  }

  //
  // Get Active Broker Account
  //
  GetActiveBrokerAccount() : BrokerAccount {
    return this.activeBrokerAccount
  }  
}

/* End File */