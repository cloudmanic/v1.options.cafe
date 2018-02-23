//
// Date: 2/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { StateService } from '../../../providers/state/state.service';

@Component({
  selector: 'app-trading-dashboard-chart',
  templateUrl: './dashboard-chart.component.html'
})

export class DashboardChartComponent implements OnInit {
  
  //
  // Constructor....
  //
  constructor(private stateService: StateService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
     
  }  

}

/* End File */