//
// Date: 9/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Order } from '../../../../../models/order';
import { TradeGroup } from '../../../../../models/trade-group';
import { Position } from '../../../../../models/position';
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

  constructor() { }

  ngOnInit() {
  }

}

/* End File */
