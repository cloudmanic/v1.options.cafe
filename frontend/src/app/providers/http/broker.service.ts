//
// Date: 2/13/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Broker } from '../../models/broker';
import { Balance } from '../../models/balance';
import { BrokerAccount } from '../../models/broker-account';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class BrokerService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get brokers
  //
  get() : Observable<Broker[]> {
    return this.http.get<Broker[]>(environment.app_server + '/api/v1/brokers').map(
      (data) => { return Broker.buildForEmit(data); 
    });
  }

  //
  // Get broker balances
  //
  getBalances() : Observable<Balance[]> {
    return this.http.get<Balance[]>(environment.app_server + '/api/v1/brokers/balances').map(
      (data) => { return Balance.buildForEmit(data); 
    });
  }  
}

/* End File */