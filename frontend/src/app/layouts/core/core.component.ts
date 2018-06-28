//
// Date: 6/27/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { StateService } from '../../providers/state/state.service';

@Component({
  selector: 'app-layouts-core',
  templateUrl: './core.component.html'
})

export class LayoutCoreComponent implements OnInit 
{
  successMsg: string = "";

  constructor(private stateService: StateService) {}

  //
  // NgInit...
  //
  ngOnInit() 
  {
    // Subscribe SiteSuccess events
    this.stateService.SiteSuccess.subscribe(data => {
      this.successMsg = data;
    });
  }

  //
  // Close Success
  //
  closeSuccess()
  {
    this.successMsg = "";
  }

}

/* End File */