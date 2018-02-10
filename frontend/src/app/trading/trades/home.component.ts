//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { TradeGroup } from '../../models/trade-group';
import { TradeGroupService } from '../../providers/http/trade-group.service';

@Component({
  selector: 'app-trades',
  templateUrl: './home.component.html'
})
export class TradesComponent implements OnInit {

  tradesList: TradeGroup[];
  
  //
  // Construct
  //
  constructor(private tradeGroupService: TradeGroupService) {}

  //
  // On Init
  //
  ngOnInit() {
    // Get tradegroup data
    this.tradeGroupService.get('open_date', 'desc').subscribe((data) => {
      this.tradesList = data
    });
  }

}

/* End File */