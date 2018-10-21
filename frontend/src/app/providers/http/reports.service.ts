//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import * as moment from 'moment';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { SummaryYearly, ProfitLoss } from '../../models/reports';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class ReportsService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Return a list of years we have trade groups for.
  //
  getTradeGroupYears(brokerAccount: number): Observable<number[]> {
    return this.http.get<number[]>(environment.app_server + '/api/v1/reports/' + brokerAccount + '/tradegroup/years')
      .map((data) => { return data; });
  } 

  //
  // Get a summary of a broker account by year.
  //
  getSummaryByYear(brokerAccount: number, year: number): Observable<SummaryYearly> 
  {
    return this.http.get<SummaryYearly>(environment.app_server + '/api/v1/reports/' + brokerAccount + '/summary/yearly/' + year)
      .map((data) => { return new SummaryYearly().fromJson(data); });
  }

  //
  // Return profit and loss data.
  //
  getProfitLoss(brokerAccount: number, start: Date, end: Date, group: string, sort: string, cumulative: boolean): Observable<ProfitLoss[]> 
  {
    let startDate = moment(start).format('YYYY-MM-DD');
    let endDate = moment(end).format('YYYY-MM-DD');

    // Set cum string.
    let cum = "false";

    if(cumulative)
    {
      cum = "true";
    }

    let getParms = "sort=" + sort + "&group=" + group + "&start=" + startDate + "&end=" + endDate + "&cumulative=" + cum;

    return this.http.get<ProfitLoss[]>(environment.app_server + '/api/v1/reports/' + brokerAccount + '/profit-loss?' + getParms)
      .map((data) => { return new ProfitLoss().fromJsonList(data); });
  }      
}

/* End File */