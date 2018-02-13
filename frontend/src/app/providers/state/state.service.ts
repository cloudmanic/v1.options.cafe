//
// Date: 2/13/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Injectable } from '@angular/core';
import { Broker } from '../../models/broker';
import { Watchlist } from '../../models/watchlist';
import { TradeGroup } from '../../models/trade-group';
import { BrokerAccount } from '../../models/broker-account';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { environment } from '../../../environments/environment';
import { AppService } from '../../providers/websocket/app.service';

@Injectable()
export class StateService  
{ 
  private quotes = {};
  private activeWatchlist: Watchlist;
  private activeBrokerAccount: BrokerAccount;
  
  // Trade Group stuff
  private tradeGroupSearchTerm : string;
  private activeTradeGroupList: TradeGroup[];

  //
  // Construct.
  //
  constructor(private appService: AppService) { 
    this.activeBrokerAccount = null;
  }

  //
  // Set a quote
  //
  SetQuote(data) {
    this.quotes[data.symbol] = data;
  }

  //
  // Get quotes
  //
  GetQuotes() {
    return this.quotes;
  }

  //
  // Set a trade group search term
  //
  SetTradeGroupSearchTerm(term: string) {
    this.tradeGroupSearchTerm = term;
  }

  //
  // Get trade group search term
  //
  GetTradeGroupSearchTerm() {
    return this.tradeGroupSearchTerm;
  }

  //
  // Get active tradegroup
  //
  GetActiveTradeGroups() : TradeGroup[] {
    return this.activeTradeGroupList;
  }

  //
  // Set active tradegroup
  //
  SetActiveTradeGroups(tg: TradeGroup[]) {
    this.activeTradeGroupList = tg;
  }

  //
  // Get active watchlist
  //
  GetActiveWatchlist() : Watchlist {
    return this.activeWatchlist;
  }

  //
  // Get active watchlist ID
  //
  GetActiveWatchlistId() {
    return localStorage.getItem('active_watchlist')
  }

  //
  // Set active watchlist
  //
  SetActiveWatchlist(watchlist: Watchlist) {
    this.activeWatchlist = watchlist;
    localStorage.setItem('active_watchlist', watchlist.Id);
  }

  //
  // Get stored active account id
  //
  GetStoredActiveAccountId() : number {
    return localStorage.getItem('active_account')
  }

  //
  // Set Active Broker Account
  //
  SetActiveBrokerAccount(brokerAccount: BrokerAccount) {
    this.activeBrokerAccount = brokerAccount
    localStorage.setItem('active_account', brokerAccount.Id);
    this.appService.RequestAllData();

    // Clear cached data.
    this.activeTradeGroupList = []
    this.tradeGroupSearchTerm = ""  
  }

  //
  // Get Active Broker Account
  //
  GetActiveBrokerAccount() : BrokerAccount {
    return this.activeBrokerAccount
  }  
}

/* End File */