//
// Date: 10/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import * as moment from 'moment';
import { Me } from '../../models/me';
import { Coupon } from '../../models/coupon';
import { Subscription } from '../../models/subscription';
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
  // Get user subscription
  //
  getSubscription(): Observable<Subscription> {
    return this.http.get<Subscription>(environment.app_server + '/api/v1/me/subscription')
      .map((data) => { return new Subscription().fromJson(data); });
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

  //
  // Change password
  //
  restPassword(current_pass: string, new_pass: string): Observable<boolean> {

    let post = {
      current_password: current_pass,
      new_password: new_pass
    }

    return this.http.put<boolean>(environment.app_server + '/api/v1/me/rest-password', post)
      .map((data) => { return true; });
  }

  //
  // Verify Coupon Code
  //
  getVerifyCoupon(code: string): Observable<Coupon> {
    return this.http.get<Coupon>(environment.app_server + '/api/v1/me/verify-coupon/' + code)
      .map((data) => {
        let coupon = new Coupon();
        coupon.Valid = data["valid"];
        coupon.Name = data["name"];
        coupon.Code = data["code"];
        coupon.AmountOff = data["amount_off"];
        coupon.PercentOff = data["percent_off"];
        coupon.Duration = data["duration"];
        return coupon; 
      });
  }

  //
  // Apply coupon.
  //
  applyCoupon(code: string): Observable<boolean> {
    let post = {
      coupon_code: code
    }

    return this.http.post<boolean>(environment.app_server + '/api/v1/me/apply-coupon', post)
      .map((data) => { return true; });
  }

  //
  // Update credit card.
  //
  updateCreditCard(token: string, coupon: string): Observable<boolean> {
    let post = {
      token: token,
      coupon_code: coupon 
    }

    return this.http.put<boolean>(environment.app_server + '/api/v1/me/update-credit-card', post)
      .map((data) => { return true; });
  }             
}

/* End File */