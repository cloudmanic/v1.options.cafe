//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
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
  limit: number = 0;
  noLimitCount: number = 0;  
  tradesList: TradeGroup[];
  searchTerm: string = ""
  tradeSelect: string = "All"
  activeAccount: BrokerAccount
  destory: Subject<boolean> = new Subject<boolean>();
  
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

    // Subscribe to changes in the selected broker.
    this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
      this.getTradeGroups();
    });

    // Load tradegroups from server
    this.getTradeGroups();
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
  // Get trade groups
  //
  getTradeGroups() 
  {
    // Get tradegroup data
    this.tradeGroupService.get(Number(this.stateService.GetStoredActiveAccountId()), 25, this.page, 'open_date', 'desc', this.searchTerm, this.tradeSelect).subscribe((res) => {
      this.limit = res.Limit;
      this.noLimitCount = res.NoLimitCount;
      this.tradesList = res.Data;
      this.count = res.Data.length;      
      this.stateService.SetActiveTradeGroups(res.Data);
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