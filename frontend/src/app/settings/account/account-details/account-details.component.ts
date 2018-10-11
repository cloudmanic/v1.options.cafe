//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Me } from '../../../models/me';
import { Subscription } from '../../../models/subscription';
import { Component, OnInit } from '@angular/core';
import { MeService } from '../../../providers/http/me.service';
import { StateService } from '../../../providers/state/state.service';

@Component({
  selector: '[app-settings-account-account-details]',
  templateUrl: './account-details.component.html',
  styleUrls: []
})

export class AccountDetailsComponent implements OnInit 
{
  hasPlan: boolean = false;
  showCloseDownAccount: boolean = false;
  userSubscription: Subscription = new Subscription(); 

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
    // Load page data.
    this.getSubscriptionData();
  }

  //
  // Get subscription data.
  //
  getSubscriptionData()
  {
    // Ajax call to get the subscription data.
    this.meService.getSubscription().subscribe((res) => {
      this.userSubscription = res;
      console.log(res);
    });
  }

  //
  // Cancel account.
  //
  cancelAccount()
  {
    this.showCloseDownAccount = true;
  }

  //
  // Close cancel account overlay.
  //
  doCancelAccount() {
    this.showCloseDownAccount = false;
  }  
}

/* End File */