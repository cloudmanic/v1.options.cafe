//
// Date: 11/10/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import * as moment from 'moment';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Symbol } from '../../models/symbol';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class SymbolService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get a symbol
  //
  getSymbol(symbol: string): Observable<Symbol> {
    return this.http.get<Symbol>(environment.app_server + '/api/v1/symbols/' + symbol).map((data) => {
      return new Symbol().fromJson(data);
    });
  }

  //
  // Search for symbols
  //
  searchSymbols(query: string) : Observable<Symbol[]> 
  {
    return this.http.get<Symbol[]>(environment.app_server + '/api/v1/symbols?search=' + query).map((data) => { 
      return new Symbol().fromJsonList(data); 
    });
  }

  //
  // Add an active symbol
  //
  addActiveSymbol(symbol: string): Observable<AddActiveSymbolResponse> 
  {
    return this.http.post<AddActiveSymbolResponse>(environment.app_server + '/api/v1/symbols/add-active-symbol', { symbol: symbol }).map((data) => {
      let r = new AddActiveSymbolResponse();
      r.Id = data["id"];
      r.Symbol = data["symbol"];
      return r;
    });
  }

  //
  // Add an active symbol
  //
  getOptionSymbolFromParts(symbol: string, expire: Date, strike: number, type: string): Observable<Symbol> 
  {
    // Format expire
    let expr = moment(new Date(expire)).format("YYYY-MM-DD"); 

    // Send AJAX Call
    return this.http.post<Symbol>(environment.app_server + '/api/v1/symbols/get-option-symbol-from-parts', { symbol: symbol, expire: expr, strike: strike, type: type }).map((data) => {
      return new Symbol().fromJson(data);
    });
  }  
}

//
// Response of a add Symbol
//
export class AddActiveSymbolResponse {
  Id: number;
  Symbol: string;
}

/* End File */