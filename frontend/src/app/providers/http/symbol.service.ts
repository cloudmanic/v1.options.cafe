//
// Date: 11/10/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
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
  // Search for symbols
  //
  searchSymbols(query: string) : Observable<Symbol[]> {
    return this.http.get<Symbol[]>(environment.app_server + '/api/v1/symbols?search=' + query).map(
      (data) => { return new Symbol().fromJsonList(data); 
    });
  }
}

/* End File */