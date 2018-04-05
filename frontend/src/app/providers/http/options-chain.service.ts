//
// Date: 4/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import * as moment from 'moment';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Symbol } from '../../models/symbol';
import { OptionsChain } from '../../models/options-chain';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class OptionsChainService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get options chain by symbol / expire date
  //
  getChainBySymbolExpire(symbol: string, expire: Date) : Observable<any> 
  {
    let expr = moment(new Date(expire)).format("YYYY-MM-DD"); 

    return this.http.get<OptionsChain>(environment.app_server + '/api/v1/quotes/options/chain/' + symbol + '/' + expr)
      .map((data) => { return new OptionsChain().fromJson(data); });
  }    
}

/* End File */