//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { BillingHistory } from '../../../models/billing-history';
import { MeService } from '../../../providers/http/me.service';
import { StateService } from '../../../providers/state/state.service';

@Component({
  selector: '[app-settings-account-billing-history]',
  templateUrl: './billing-history.component.html',
  styleUrls: []
})
export class BillingHistoryComponent implements OnInit 
{
  billingHistory: BillingHistory[] = [];

  //
  // Construct.
  //
  constructor(private stateService: StateService, private meService: MeService) {
    // Get cached data
    //this.userProfile = this.stateService.GetSettingsUserProfile();
  }

  // 
  // NG Init.
  //
  ngOnInit() 
  {
    this.getInvoiceData();
  }

  //
  // Get invoice data.
  //
  getInvoiceData()
  {
    // Ajax call to get the billing history data.
    this.meService.getBillingHistory().subscribe((res) => {
      this.billingHistory = res;
    });
  }

}

/* End File */