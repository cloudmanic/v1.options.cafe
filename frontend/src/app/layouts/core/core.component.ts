//
// Date: 6/27/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { StateService } from '../../providers/state/state.service';
import { NotificationsService } from '../../providers/http/notifications.service';

@Component({
  selector: 'app-layouts-core',
  templateUrl: './core.component.html'
})

export class LayoutCoreComponent implements OnInit 
{
  successMsg: string = "";

  //
  // Construct
  //
  constructor(private stateService: StateService, private notificationsService: NotificationsService) { }

  //
  // NgInit...
  //
  ngOnInit() 
  {
    // Subscribe SiteSuccess events
    this.stateService.SiteSuccess.subscribe(data => {
      this.successMsg = data;
      setTimeout(() => { this.closeSuccess() }, 3000);
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