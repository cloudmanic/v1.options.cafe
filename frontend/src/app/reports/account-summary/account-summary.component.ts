//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Rx';
import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit } from '@angular/core';
import { SummaryYearly } from '../../models/reports';
import { StateService } from '../../providers/state/state.service';
import { ReportsService } from '../../providers/http/reports.service';
import { DropdownAction } from '../../shared/dropdown-select/dropdown-select.component';

@Component({
  selector: 'app-account-summary',
  templateUrl: './account-summary.component.html',
  styleUrls: []
})

export class AccountSummaryComponent implements OnInit 
{
  tradeGroupYears: number[];
  summaryByYear: SummaryYearly;
  summaryByYearSelected: number;
  summaryActions: DropdownAction[] = null;

  private destory: Subject<boolean> = new Subject<boolean>(); 

  //
  // Construct.
  //
  constructor(private stateService: StateService, private reportsService: ReportsService) 
  {
    // Get data from site state.
    this.summaryByYear = this.stateService.GetReportsSummaryByYear();
    this.summaryByYearSelected = this.stateService.GetReportsSummaryByYearSelectedYear();
  }

  //
  // NG Init
  //
  ngOnInit() 
  {
    // Subscribe to changes in the selected broker.
    this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
      this.getAccountSummary();
      this.getTradeGroupYears();
    });

    // Get data on page load.
    this.getAccountSummary();
    this.getTradeGroupYears();
  }

  //
  // OnDestroy
  //
  ngOnDestroy() {
    this.destory.next();
    this.destory.complete();
  } 

  //
  // Setup Summary actions.
  //
  setupSummaryActions() {
    let das = []
    this.summaryActions = []

    // Loop through add dates to the drop down
    for(let i = 0; i < this.tradeGroupYears.length; i++)
    {
      // First action
      let da = new DropdownAction();
      da.title = 'Year ' + this.tradeGroupYears[i];

      // Click on year
      da.click = (row: number[]) => {
        this.summaryByYearSelected = this.tradeGroupYears[i];
        this.getAccountSummary();
        this.stateService.SetReportsSummaryByYearSelectedYear(this.summaryByYearSelected);

        // Hack to get it to close;
        this.setupSummaryActions();
      };

      das.push(da);
    }

    this.summaryActions = das;
  }  

  //
  // Get Data = Trade Group Years
  //
  getTradeGroupYears()
  {
    // Make api call to get years
    this.reportsService.getTradeGroupYears(Number(this.stateService.GetStoredActiveAccountId())).subscribe((res) => {
      this.tradeGroupYears = res;
      this.setupSummaryActions();
    });    
  }

  //
  // Get Data - Account Summary
  //
  getAccountSummary() 
  {
    // Make api call to get account summary
    this.reportsService.getSummaryByYear(Number(this.stateService.GetStoredActiveAccountId()), this.summaryByYearSelected).subscribe((res) => {
      this.summaryByYear = res;
      this.stateService.SetReportsSummaryByYear(this.summaryByYear);
    });
  }

}

/* End File */
