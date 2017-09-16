//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AppService } from '../../providers/websocket/app.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './home.component.html'
})
export class DashboardComponent implements OnInit {

  constructor(private app: AppService) { }

  ngOnInit() {
    
  }

}
