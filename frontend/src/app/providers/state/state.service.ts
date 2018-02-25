//
// Date: 2/13/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Injectable } from '@angular/core';
import { Broker } from '../../models/broker';
import { Symbol } from '../../models/symbol';
import { Watchlist } from '../../models/watchlist';
import { HistoricalQuote } from '../../models/historical-quote';
import { TradeGroup, TradeGroupsCont } from '../../models/trade-group';
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

  // Dashboard stuff
  private dashboardChartRangeSelect: string = "days-monthly-365";
  private dashboardChartData: HistoricalQuote[];
  private dashboardTradeGroups: TradeGroupsCont;
  private dashboardChartSymbol: Symbol = new Symbol(1, "SPDR S&P 500 ETF Trust", "SPY", "Equity", null);
  
  // Trade Group stuff
  private tradeGroupPage: number = 1;
  private tradeGroupSearchTerm: string = "";
  private tradeGroupTradeSelect: string = "All";
  private activeTradeGroupList: TradeGroup[];

  // Emitters - Pushers
  public BrokerChange = new EventEmitter<number>();

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
  // Set a dashboard chart RangeSelect
  //
  SetDashboardChartRangeSelect(data: string) {
    this.dashboardChartRangeSelect = data;
  }

  //
  // Get dashboard chart RangeSelect
  //
  GetDashboardChartRangeSelect() {
    return this.dashboardChartRangeSelect;
  }

  //
  // Set a dashboard chart symbol
  //
  SetDashboardChartSymbol(data: Symbol) {
    this.dashboardChartSymbol = data;
  }

  //
  // Get dashboard chart symbol
  //
  GetDashboardChartSymbol() {
    return this.dashboardChartSymbol;
  }

  //
  // Set dashboard chart data
  //
  SetDashboardChartData(data: HistoricalQuote[]) {
    this.dashboardChartData = data;
  }

  //
  // Get dashboard chart data
  //
  GetDashboardChartData() {
    return this.dashboardChartData;
  }

  //
  // Set a trade group search term
  //
  SetDashboardTradeGroups(group: TradeGroupsCont) {
    this.dashboardTradeGroups = group;
  }

  //
  // Get trade group page
  //
  GetDashboardTradeGroups() {
    return this.dashboardTradeGroups;
  }

  //
  // Set a trade group search term
  //
  SetTradeGroupSearchTerm(term: string) {
    this.tradeGroupSearchTerm = term;
  }

  //
  // Get trade group page
  //
  GetTradeGroupPage() {
    return this.tradeGroupPage;
  }

  //
  // Set a trade group page
  //
  SetTradeGroupPage(page: number) {
    this.tradeGroupPage = page;
  }

  //
  // Get trade group search term
  //
  GetTradeGroupSearchTerm() {
    return this.tradeGroupSearchTerm;
  }

  //
  // Set a trade group trade select
  //
  SetTradeGroupTradeSelect(option: string) {
    this.tradeGroupTradeSelect = option;
  }

  //
  // Get trade group trade select
  //
  GetTradeGroupTradeSelect() {
    return this.tradeGroupTradeSelect;
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
    localStorage.setItem('active_watchlist', String(watchlist.Id));
  }

  //
  // Get stored active account id
  //
  GetStoredActiveAccountId() : string {
    return localStorage.getItem('active_account');
  }

  //
  // Set Active Broker Account
  //
  SetActiveBrokerAccount(brokerAccount: BrokerAccount) {
    this.activeBrokerAccount = brokerAccount
    localStorage.setItem('active_account', String(brokerAccount.Id));
    
    this.BrokerChange.emit(brokerAccount.Id);

    // Clear cached data.
    this.activeTradeGroupList = [];
    this.tradeGroupSearchTerm = "";
    this.tradeGroupTradeSelect = "All";  
  }

  //
  // Get Active Broker Account
  //
  GetActiveBrokerAccount() : BrokerAccount {
    return this.activeBrokerAccount
  }  
}

/* End File */