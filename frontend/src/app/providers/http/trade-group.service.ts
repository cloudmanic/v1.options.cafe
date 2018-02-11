//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { TradeGroup } from '../../models/trade-group';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class TradeGroupService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get trade groups
  //
  get(page: number, order: string, sort: string, search: string) : Observable<TradeGroup[]> {
    return this.http.get<TradeGroup[]>(environment.app_server + '/api/v1/tradegroups?page=' + page + '&order=' + order + '&sort=' + sort + '&search=' + search).map(
      (data) => { return TradeGroup.buildForEmit(data); 
    });
  }
}

/* End File */