//
// Date: 2/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit } from '@angular/core';
import { Order } from '../../../models/order';
import { ChangeDetected } from '../../../models/change-detected';
import { TradeGroup, TradeGroupsCont } from '../../../models/trade-group';
import { WebsocketService } from '../../../providers/http/websocket.service';
import { StateService } from '../../../providers/state/state.service';
import { TradeGroupService } from '../../../providers/http/trade-group.service';

@Component({
  selector: 'app-trading-positions',
  templateUrl: './positions.component.html'
})

export class PositionsComponent implements OnInit {
  
  public quotes = {}
  public orders: Order[];  
  public tradeGroups: TradeGroupsCont;

  private destory: Subject<boolean> = new Subject<boolean>();

  //
  // Constructor....
  //
  constructor(private websocketService: WebsocketService, private stateService: StateService, private tradeGroupService: TradeGroupService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Get the Positions
    this.getPositions();

    // Get Data from cache
    this.quotes = this.stateService.GetQuotes();    
    this.tradeGroups = this.stateService.GetDashboardTradeGroups();

    // Subscribe to changes in the selected broker.
    this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
      this.getPositions();
    });

    // Subscribe to when changes are detected at the server.
    this.websocketService.changedDetectedPush.takeUntil(this.destory).subscribe(data => {
      this.manageChangeDetection(data);
    }); 

    // Subscribe to data updates from the quotes - Market Quotes
    this.websocketService.quotePushData.takeUntil(this.destory).subscribe(data => {
      this.quotes[data.symbol] = data;
    });     
  }

  //
  // OnDestroy
  //
  ngOnDestroy()
  {
    this.destory.next();
    this.destory.complete();
  }

  //
  // Manage change detection.
  //
  private manageChangeDetection(data: ChangeDetected)
  {
    if(data.Type == "orders") 
    {
      this.getPositions();
    }
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

      // Store the tradegroups in the state manager
      this.stateService.SetDashboardTradeGroups(this.tradeGroups);
    });    
  }
}

/* End File */