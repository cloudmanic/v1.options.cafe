//
// Date: 2/26/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { MarketStatus } from '../../models/market-status';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class StatusService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get brokers
  //
  getMarketStatus() : Observable<MarketStatus> {
    return this.http.get<MarketStatus>(environment.app_server + '/api/v1/status/market').map(
      (data) => { return MarketStatus.buildForEmit(data); 
    });
  } 
}

/* End File */