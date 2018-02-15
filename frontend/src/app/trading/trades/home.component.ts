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
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-trades',
  templateUrl: './home.component.html'
})
export class TradesComponent implements OnInit {

  page: number = 1;
  count: number = 0;
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
    // Set the page
    this.page = this.stateService.GetTradeGroupPage();

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
    this.tradeGroupService.get(Number(this.stateService.GetStoredActiveAccountId()), this.page, 'open_date', 'desc', this.searchTerm, this.tradeSelect).subscribe((data) => {
      this.tradesList = data;
      this.count = data.length;      
      this.stateService.SetActiveTradeGroups(data);
      this.stateService.SetTradeGroupPage(this.page);
    });    
  }

  //
  // On Trade select...
  //
  onTradeSelect(event) {
    this.page = 1;
    this.getTradeGroups();
    this.stateService.SetTradeGroupTradeSelect(this.tradeSelect);    
  }

  //
  // On search...
  //
  onSearchKeyUp(event) {
    this.page = 1;    
    this.getTradeGroups();
    this.stateService.SetTradeGroupSearchTerm(this.searchTerm);
  }

  //
  // On paging click.
  //
  onPagingClick(page: number) 
  {
    this.page = page;
    this.getTradeGroups();
  }   
}

/* End File */