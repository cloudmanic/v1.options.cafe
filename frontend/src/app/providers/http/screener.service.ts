//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import * as moment from 'moment';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Screener } from '../../models/screener';
import { ScreenerResult } from '../../models/screener-result';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class ScreenerService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get a list of screens in the system. 
  //
  get(): Observable<Screener[]> 
  {
    return this.http.get<Screener[]>(environment.app_server + '/api/v1/screeners')
      .map((data) => { return new Screener().fromJsonList(data); });
  } 

  //
  // Get screener results.
  //
  getResults(id: number): Observable<ScreenerResult[]> {
    return this.http.get<ScreenerResult[]>(environment.app_server + '/api/v1/screeners/' + id + '/results')
      .map((data) => { return new ScreenerResult().fromJsonList(data); });
  }     
}

/* End File */