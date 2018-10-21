//
// Date: 10/20/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { ProfitLoss } from '../../models/reports';
import { Component, OnInit } from '@angular/core';
import { StateService } from '../../providers/state/state.service';
import { ReportsService } from '../../providers/http/reports.service';

@Component({
  selector: 'app-reports-custom-reports',
  templateUrl: './custom-reports.component.html',
  styleUrls: []
})

export class CustomReportsComponent implements OnInit 
{
  listData: ProfitLoss[] = [];

  //
  // Construct.
  //
  constructor(private stateService: StateService, private reportsService: ReportsService) 
  {

  }

  //
  // NG Init
  //
  ngOnInit() 
  {
    // Get data for page.
    this.getProfitLoss();
  }

  //
  // Get Data = Profit Loss
  //
  getProfitLoss() 
  {
    this.reportsService.getProfitLoss(Number(this.stateService.GetStoredActiveAccountId()), "2018-01-01", "2018-12-31", "month", "desc").subscribe((res) => {
      this.listData = res;

      console.log(res);
    });
  }

  //
  // Get profit total
  //
  getProfitTotal(rows: ProfitLoss[]): number 
  {
    let total = 0;

    for(let i = 0; i < rows.length; i++)
    {
      total += rows[i].Profit;
    }

    return total;
  }

  //
  // Get trade total
  //
  getTradeTotal(rows: ProfitLoss[]): number 
  {
    let total = 0;

    for (let i = 0; i < rows.length; i++) 
    {
      total += rows[i].TradeCount;
    }

    return total;
  }

  //
  // Get trade total
  //
  getCommissionsTotal(rows: ProfitLoss[]): number 
  {
    let total = 0;

    for (let i = 0; i < rows.length; i++) 
    {
      total += rows[i].Commissions;
    }

    return total;
  }

}

/* End File */