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

@Component({
  selector: 'app-account-summary',
  templateUrl: './account-summary.component.html',
  styleUrls: []
})

export class AccountSummaryComponent implements OnInit 
{
  summaryByYear: SummaryYearly;
  summaryByYearSelected: number;

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
      this.stateService.GetReportsSummaryByYearSelectedYear()
    });

    // Get data on page load.
    this.getAccountSummary();
  }

  //
  // OnDestroy
  //
  ngOnDestroy() {
    this.destory.next();
    this.destory.complete();
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
