//
// Date: 2/13/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Order } from '../../models/order';
import { Broker } from '../../models/broker';
import { Balance } from '../../models/balance';
import { BrokerAccount } from '../../models/broker-account';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class BrokerService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get brokers
  //
  get() : Observable<Broker[]> {
    return this.http.get<Broker[]>(environment.app_server + '/api/v1/brokers').map(
      (data) => { return Broker.buildForEmit(data); 
    });
  }

  //
  // Get broker balances
  //
  getBalances(brokerId: number) : Observable<Balance[]> {
    return this.http.get<Balance[]>(environment.app_server + '/api/v1/brokers/' + brokerId + '/balances').map(
      (data) => { return Balance.buildForEmit(data); 
    });
  }

  //
  // Get broker account balance
  //
  getAccountBalance(brokerId: number, brokerAccountId: number): Observable<Balance> {
    return this.http.get<Balance>(environment.app_server + '/api/v1/brokers/' + brokerId + '/accounts/' + brokerAccountId + '/balance').map(
      (data) => {
        return Balance.fromJson(data);
      });
  }

  //
  // Get broker orders
  //
  getOrders(brokerId: number) : Observable<Order[]> {
    return this.http.get<Order[]>(environment.app_server + '/api/v1/brokers/' + brokerId + '/orders').map(
      (data) => { return Order.buildForEmit(data); 
    });
  }

  //
  // Create a new broker
  //
  create(name: string, display_name: string): Observable<Broker> {
    let body = {
      name: name,
      display_name: display_name
    }

    return this.http.post<Broker>(environment.app_server + '/api/v1/brokers', body)
      .map((data) => { return Broker.fromJson(data); });
  }  

  //
  // Update a broker
  //
  update(brokerId: number, display_name: string): Observable<Broker> {
    let body = {
      display_name: display_name
    }

    return this.http.put<Broker>(environment.app_server + '/api/v1/brokers/' + brokerId, body)
      .map((data) => { return Broker.fromJson(data); });
  }

  //
  // Update a broker account
  //
  updateBrokerAccount(brokerId: number, brokerAccountId: number, name: string, stock_commission: number, stock_min: number, option_commission: number, option_single_min: number, option_multi_leg_min: number, option_base: number): Observable<Broker> {
    let body = {
      name: name, 
      stock_commission: stock_commission, 
      stock_min: stock_min, 
      option_commission: option_commission, 
      option_single_min: option_single_min, 
      option_multi_leg_min: option_multi_leg_min, 
      option_base: option_base 
    }

    return this.http.put<Broker>(environment.app_server + '/api/v1/brokers/' + brokerId + '/accounts/' + brokerAccountId, body)
      .map((data) => { return Broker.fromJson(data); });
  }       
}

/* End File */