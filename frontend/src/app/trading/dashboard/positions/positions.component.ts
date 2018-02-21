//
// Date: 2/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { Order } from '../../../models/order';
import { TradeGroup } from '../../../models/trade-group';
import { AppService } from '../../../providers/websocket/app.service';
import { QuoteService } from '../../../providers/websocket/quote.service';
import { StateService } from '../../../providers/state/state.service';
import { TradeGroupService } from '../../../providers/http/trade-group.service';

@Component({
  selector: 'app-trading-positions',
  templateUrl: './positions.component.html'
})

export class PositionsComponent implements OnInit {
  
  private quotes = {}
  private orders: Order[];  
  private tradeGroups: TradeGroupsCont;

  //
  // Constructor....
  //
  constructor(private appService: AppService, private quoteService: QuoteService, private stateService: StateService, private tradeGroupService: TradeGroupService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Get the Positions
    this.getPositions();

    // // Get Data from cache
    // this.setOrders(this.appService.orders);
    // this.quotes = this.quoteService.quotes;
        
    // Subscribe to data updates from the broker - Orders
    this.appService.ordersPush.subscribe(data => {
      this.orders = data;
    });    
    
    // Subscribe to data updates from the quotes - Market Quotes
    this.quoteService.marketQuotePushData.subscribe(data => {
      this.quotes[data.symbol] = data;
    });     
  }

  //
  // Get positions AKA Trade Groups
  //
  private getPositions() 
  {
    // Get tradegroup data
    this.tradeGroupService.get(Number(this.stateService.GetStoredActiveAccountId()), 100, 1, 'open_date', 'asc', '', 'Open').subscribe((res) => {

      // Reset the trade groups
      this.tradeGroups = new TradeGroupsCont([], [], [], [], [], [], []);

      // Loop through and classify positions
      for(let i = 0; i <= res.Data.length; i++)
      {
        // Not sure why I have to do this.
        if(typeof res.Data[i] == "undefined")
        {
          continue;
        }

        // Push onto the array.
        this.tradeGroups[res.Data[i].Type.split(' ').join('')].push(res.Data[i]);

      }

      // this.limit = res.Limit;
      // this.noLimitCount = res.NoLimitCount;
      // this.tradesList = res.Data;
      // this.count = res.Data.length;      
      // this.stateService.SetActiveTradeGroups(res.Data);
      // this.stateService.SetTradeGroupPage(this.page);
    });    
  }
}

//
// Setup a class to hold all the different position types
//
export class TradeGroupsCont 
{
  constructor(
    public Option: TradeGroup[],
    public PutCreditSpread: TradeGroup[],
    public CallCreditSpread: TradeGroup[], 
    public PutDebitSpread: TradeGroup[], 
    public CallDebitSpread: TradeGroup[], 
    public IronCondor: TradeGroup[], 
    public Other: TradeGroup[],         
  ){}  
}

/* End File */