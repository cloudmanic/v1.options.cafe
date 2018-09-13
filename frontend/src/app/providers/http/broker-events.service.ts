//
// Date: 9/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import * as moment from 'moment';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { BrokerEvent } from '../../models/broker-event';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class BrokerEventsService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }       

  //
  // Get broker events
  //
  get(broker_account_id: number, limit: number, page: number, order: string, sort: string, search: string): Observable<BrokerEventsResponse> {

    // Make API call.
    return this.http.get<Object[]>(environment.app_server + '/api/v1/broker-events/' + broker_account_id + '?limit=' + limit + '&page=' + page + '&order=' + order + '&sort=' + sort, { observe: 'response' }).map((res) => {
      let lastPage = false;

      // Build last page
      if (res.headers.get('X-Last-Page') == "true") 
      {
        lastPage = true;
      }

      // Build and return data
      return new BrokerEventsResponse(lastPage, Number(res.headers.get('X-Offset')), Number(res.headers.get('X-Limit')), Number(res.headers.get('X-No-Limit-Count')), new BrokerEvent().fromJsonList(res.body));
    });
  }
}

//
// Broker Events Response
//
export class BrokerEventsResponse {
  constructor(
    public LastPage: boolean,
    public Offset: number,
    public Limit: number,
    public NoLimitCount: number,
    public Data: BrokerEvent[]
  ) { }
}

/* End File */