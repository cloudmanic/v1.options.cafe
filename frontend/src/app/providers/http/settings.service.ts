//
// Date: 10/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import * as moment from 'moment';
import { Settings } from '../../models/settings';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class SettingsService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get settings
  //
  get(): Observable<Settings> 
  {
    return this.http.get<Settings>(environment.app_server + '/api/v1/settings')
      .map((data) => { return new Settings().fromJson(data); });
  }             
}

/* End File */