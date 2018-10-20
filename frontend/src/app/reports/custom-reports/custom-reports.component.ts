//
// Date: 10/20/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { SummaryYearly } from '../../models/reports';
import { StateService } from '../../providers/state/state.service';
import { ReportsService } from '../../providers/http/reports.service';

@Component({
  selector: 'app-reports-custom-reports',
  templateUrl: './custom-reports.component.html',
  styleUrls: []
})

export class CustomReportsComponent implements OnInit 
{
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

  }

}

/* End File */