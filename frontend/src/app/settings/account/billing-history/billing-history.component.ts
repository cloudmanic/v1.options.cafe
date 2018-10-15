//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { Invoice } from '../../../models/invoice';
import { MeService } from '../../../providers/http/me.service';
import { StateService } from '../../../providers/state/state.service';

@Component({
  selector: '[app-settings-account-billing-history]',
  templateUrl: './billing-history.component.html',
  styleUrls: []
})
export class BillingHistoryComponent implements OnInit 
{
  invoices: Invoice[] = [];

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
    // Ajax call to get the subscription data.
    this.meService.getInvoices().subscribe((res) => {
      this.invoices = res;
    });
  }

}

/* End File */