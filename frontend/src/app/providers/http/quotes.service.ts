//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HistoricalQuote } from '../../models/historical-quote';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class QuotesService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get historical quotes
  //
  getHistoricalQuote(symbol: string, start: Date, end: Date, interval: string) : Observable<HistoricalQuote[]> {

    let ts = "";

    let request = environment.app_server + '/api/v1/quotes/historical?symbol=' + symbol + '&start=' + start.toISOString().substring(0, 10) + 
                  '&end=' + end.toISOString().substring(0, 10) + '&interval=' + interval;

    // Make API call.
    return this.http.get<HistoricalQuote[]>(request).map(
      (data) => { return HistoricalQuote.buildForEmit(data); 
    });
  }
}

/* End File */