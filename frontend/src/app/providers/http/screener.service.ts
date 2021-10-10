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
  // Get screener by id.
  //
  getById(id: number): Observable<Screener> {
    return this.http.get<Screener>(environment.app_server + '/api/v1/screeners/' + id)
      .map((data) => { return new Screener().fromJson(data); });
  } 

  //
  // Delete screener by id.
  //
  deleteById(id: number): Observable<boolean> {
    return this.http.delete<boolean>(environment.app_server + '/api/v1/screeners/' + id).map((data) => { return true });    
  }

  //
  // Get screener results.
  //
  getResults(id: number, brokerAccount: number): Observable<ScreenerResult[]> {
    return this.http.get<ScreenerResult[]>(environment.app_server + '/api/v1/screeners/' + id + '/results?broker_account_id=' + brokerAccount)
      .map((data) => { return new ScreenerResult().fromJsonList(data); });
  }

  //
  // Submit screen to be saved.
  //
  submitScreen(screen: Screener): Observable<Screener> {

    let body = {
      name: screen.Name,
      strategy: screen.Strategy,
      symbol: screen.Symbol,
      items: []
    }

    for (let i = 0; i < screen.Items.length; i++) 
    {
      body.items.push({
        key: screen.Items[i].Settings.Key,
        operator: screen.Items[i].Operator,
        value_number: screen.Items[i].ValueNumber,
        value_string: screen.Items[i].ValueString
      });
    }

    return this.http.post<Screener>(environment.app_server + '/api/v1/screeners', body)
      .map((data) => { return new Screener().fromJson(data); });
  }

  //
  // Update screener.
  //
  submitUpdate(screen: Screener): Observable<boolean> {
    let body = {
      id: screen.Id,
      name: screen.Name,
      strategy: screen.Strategy,
      symbol: screen.Symbol,
      items: []
    }

    for (let i = 0; i < screen.Items.length; i++) 
    {
      body.items.push({
        key: screen.Items[i].Settings.Key,
        operator: screen.Items[i].Operator,
        value_number: screen.Items[i].ValueNumber,
        value_string: screen.Items[i].ValueString
      });
    }

    return this.http.put<boolean>(environment.app_server + '/api/v1/screeners/' + screen.Id, body).map((data) => { return true });
  } 

  //
  // Submit screen not saved in system.
  //
  submitScreenForResults(screen: Screener): Observable<ScreenerResult[]> {

    let body = {
      name: 'One Time',
      strategy: screen.Strategy,
      symbol: screen.Symbol,
      items: []
    }

    for (let i = 0; i < screen.Items.length; i++) 
    {
      body.items.push({
        key: screen.Items[i].Settings.Key,
        operator: screen.Items[i].Operator,
        value_number: screen.Items[i].ValueNumber,
        value_string: screen.Items[i].ValueString
      });
    }

    return this.http.post<ScreenerResult[]>(environment.app_server + '/api/v1/screeners/results', body)
      .map((data) => { return new ScreenerResult().fromJsonList(data); });
  }       
}

/* End File */