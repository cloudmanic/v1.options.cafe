//
// Date: 10/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import * as moment from 'moment';
import { Me } from '../../models/me';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class MeService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get user profile
  //
  getProfile(): Observable<Me> 
  {
    return this.http.get<Me>(environment.app_server + '/api/v1/me/profile')
      .map((data) => { return new Me().fromJson(data); });
  } 

  //
  // Save user profile
  //
  saveProfile(profile: Me): Observable<Me> {

    let post = {
      first_name: profile.FirstName,
      last_name: profile.LastName,
      email: profile.Email,
      phone: profile.Phone,
      address: profile.Address,
      city: profile.City,
      state: profile.State,
      zip: profile.Zip,
      country: profile.Country
    }

    return this.http.put<Me>(environment.app_server + '/api/v1/me/profile', post)
      .map((data) => { return new Me().fromJson(data); });
  }         
}

/* End File */