//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { TradeGroup } from '../../models/trade-group';
import { BrokerAccount } from '../../models/broker-account';
import { AppService } from '../../providers/websocket/app.service';
import { TradeGroupService } from '../../providers/http/trade-group.service';
import { StateService } from '../../providers/state/state.service';

@Component({
  selector: 'app-trades',
  templateUrl: './home.component.html'
})
export class TradesComponent implements OnInit {

  tradesList: TradeGroup[];
  searchTerm: string = ""
  tradeSelect: string = "All"
  activeAccount: BrokerAccount
  
  //
  // Construct
  //
  constructor(private appService: AppService, private tradeGroupService: TradeGroupService, private stateService: StateService) {}

  //
  // On Init
  //
  ngOnInit() 
  {
    // Set the search term from cache
    this.searchTerm = this.stateService.GetTradeGroupSearchTerm();

    // Set the cached trade select
    this.tradeSelect = this.stateService.GetTradeGroupTradeSelect();

    // Load trade groups from cache.
    this.tradesList = this.stateService.GetActiveTradeGroups(); 

    // Load tradegroups from server
    this.getTradeGroups();
  }

  //
  // Get trade groups
  //
  getTradeGroups() 
  {
    // Get tradegroup data
    this.tradeGroupService.get(Number(this.stateService.GetStoredActiveAccountId()), 1, 'open_date', 'desc', this.searchTerm, this.tradeSelect).subscribe((data) => {
      this.tradesList = data;
      this.stateService.SetActiveTradeGroups(data);
    });    
  }

  //
  // On Trade select...
  //
  onTradeSelect(event) {
    this.getTradeGroups();
    this.stateService.SetTradeGroupTradeSelect(this.tradeSelect);    
  }

  //
  // On search...
  //
  onSearchKeyUp(event) {
    this.getTradeGroups();
    this.stateService.SetTradeGroupSearchTerm(this.searchTerm);
  }  
}

/* End File */