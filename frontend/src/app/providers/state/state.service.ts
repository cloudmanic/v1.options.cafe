//
// Date: 2/13/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Rank } from '../../models/rank';
import { Order } from '../../models/order';
import { Broker } from '../../models/broker';
import { Symbol } from '../../models/symbol';
import { Screener } from '../../models/screener';
import { Watchlist } from '../../models/watchlist';
import { BrokerEvent } from '../../models/broker-event';
import { SummaryYearly } from '../../models/reports';
import { HistoricalQuote } from '../../models/historical-quote';
import { TradeGroup, TradeGroupsCont } from '../../models/trade-group';
import { BrokerAccount } from '../../models/broker-account';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { environment } from '../../../environments/environment';

@Injectable()
export class StateService  
{ 
  private quotes = {};
  private activeWatchlist: Watchlist;
  private activeBrokerAccount: BrokerAccount = null;

  // Dashboard stuff
  private dashboardVixRank: Rank;
  private dashboardChartRangeSelect: string = "today-1min-1";
  private dashboardChartData: HistoricalQuote[];
  private dashboardTradeGroups: TradeGroupsCont;
  private dashboardChartSymbol: Symbol = new Symbol().New(1, "SPDR S&P 500 ETF Trust", "SPY", "Equity", null, null, null, null);

  // Trade Group stuff
  private tradeGroupPage: number = 1;
  private tradeGroupSearchTerm: string = "";
  private tradeGroupTradeSelect: string = "All";
  private activeTradeGroupList: TradeGroup[];

  // Order stuff
  private activeOrders: Order[];

  // Screener
  private screenerScreens: Screener[];

  // Reports
  private reportsSummaryByYear: SummaryYearly;
  private reportsSummaryByYearSelected: number = new Date().getFullYear();

  // Account History
  private accountHistoryPage: number = 1;
  private accountHistoryList: BrokerEvent[];

  // Emitters - Pushers
  public SiteSuccess = new EventEmitter<string>();
  public BrokerChange = new EventEmitter<number>();


  //
  // Constructor.
  //
  constructor(private http: HttpClient, private router: Router) 
  {
    this.PingServer();
  } 

  //
  // Ping to make sure our access token is still good.
  // If not redirect back to login. Also the server
  // uses this as an opportunity to collect some stats.
  //
  PingServer()
  {

    // Ajax call to ping server every 10 seconds
    setInterval(() => {

      this.http.get(environment.app_server + '/api/v1/ping').subscribe(
        // Success
        data => {

          if(data.status == "logout")
          {
            this.router.navigate(['/logout']);
          }

        },

        // Error
        (err: HttpErrorResponse) => {

          if (err.error instanceof Error) 
          {
            // A client-side or network error occurred. Handle it accordingly.
            console.log('An error occurred:', err.error);
          } else 
          {
            // Log error
            console.log(err.error.error);

            // Access token mostly not good. 
            // If the error is blank it often means the 
            // server is down.
            if(err.error.error && (err.error.error.length > 0)) 
            {
              this.router.navigate(['/logout']);
            }
          }

        }

      );

    }, 10000);

  }

  //
  // Set accountHistoryList
  //
  SetAccountHistoryList(data: BrokerEvent[]) {
    this.accountHistoryList = data;
  }

  //
  // Get accountHistoryList
  //
  GetAccountHistoryList() {
    return this.accountHistoryList;
  }

  //
  // Set accountHistoryList Page
  //
  SetAccountHistoryPage(data: number) {
    this.accountHistoryPage = data;
  }

  //
  // Get accountHistoryList Page
  //
  GetAccountHistoryPage() {
    return this.accountHistoryPage;
  }

  //
  // Set ReportsSummaryByYear
  //
  SetReportsSummaryByYear(data: SummaryYearly) {
    this.reportsSummaryByYear = data;
  }

  //
  // Get ReportsSummaryByYear
  //
  GetReportsSummaryByYear() {
    return this.reportsSummaryByYear;
  }

  //
  // Set ReportsSummaryByYear - selected year
  //
  SetReportsSummaryByYearSelectedYear(data: number) {
    this.reportsSummaryByYearSelected = data;
  }

  //
  // Get ReportsSummaryByYear - selected year
  //
  GetReportsSummaryByYearSelectedYear() {
    return this.reportsSummaryByYearSelected;
  }

  //
  // Set Orders
  //
  SetActiveOrders(data: Order[]) {
    this.activeOrders = data;
  }

  //
  // Get Orders
  //
  GetActiveOrders() {
    return this.activeOrders;
  }

  //
  // Set Screens
  //
  SetScreens(data: Screener[]) {
    this.screenerScreens = data;
  }

  //
  // Get Screens
  //
  GetScreens() {
    return this.screenerScreens;
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
  // Set a dashboard rank
  //
  SetDashboardRank(data: Rank) {
    this.dashboardVixRank = data;
  }

  //
  // Get dashboard rank
  //
  GetDashboardRank(): Rank {
    return this.dashboardVixRank;
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
    return Number(localStorage.getItem('active_watchlist'));
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