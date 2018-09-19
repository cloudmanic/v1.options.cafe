//
// Date: 9/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Order } from '../../../../../models/order';
import { TradeGroup } from '../../../../../models/trade-group';
import { Position } from '../../../../../models/position';
import { WebsocketService } from '../../../../../providers/http/websocket.service';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../../../../../providers/http/trade.service';
import { DropdownAction } from '../../../../../shared/dropdown-select/dropdown-select.component';
import { Component, OnInit, Input, Output, OnChanges, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-trading-position-type-equity',
  templateUrl: './equity.component.html',
  styleUrls: []
})

export class EquityComponent implements OnInit {

  @Input() quotes = {};
  @Input() orders: Order[];
  @Input() tradeGroups: TradeGroup[]; 

  //
  // Constructor....
  //
  constructor(private websocketService: WebsocketService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // // Subscribe to data updates from the quotes - Market Quotes
    // this.websocketService.quotePushData.subscribe(data => {
    //   this.quotes[data.symbol] = data;
    // }); 
  }

}

/* End File */
