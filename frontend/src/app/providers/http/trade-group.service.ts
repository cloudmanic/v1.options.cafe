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
  get(broker_account_id: number, limit: number, page: number, order: string, sort: string, search: string, tradeSelect: string) : Observable<TradeGroupsResponse> {

    let ts = "";

    // Figure out the trade select url parms
    switch(tradeSelect)
    {
      case "All":
        ts = "";
      break; 

      case "Open":
        ts = "&status=Open";
      break; 
        
      case "Closed":
        ts = "&status=Closed";
      break; 

      default:
        ts = "&type=" + tradeSelect;
      break; 
    }

    // Make API call.
    return this.http.get(environment.app_server + '/api/v1/tradegroups?broker_account_id=' + broker_account_id + '&limit=' + limit + '&page=' + page + '&order=' + order + '&sort=' + sort + '&search=' + search + ts, { observe: 'response' }).map((res) => {
      let lastPage = false;

      // Build last page
      if(res.headers.get('X-Last-Page') == "true")
      {
        lastPage = true;
      }      

      // Build and return data
      return new TradeGroupsResponse(lastPage, Number(res.headers.get('X-Offset')), Number(res.headers.get('X-Limit')), Number(res.headers.get('X-No-Limit-Count')), TradeGroup.buildForEmit(res.body));
    });
  }
}

//
// Trade Groups Response
//
export class TradeGroupsResponse 
{
  constructor(
    public LastPage: boolean,
    public Offset: number,
    public Limit: number, 
    public NoLimitCount: number, 
    public Data: TradeGroup[] 
  ){}  
}

/* End File */