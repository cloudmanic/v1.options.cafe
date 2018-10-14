//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { ActivatedRoute, Params } from '@angular/router';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-settings-account-upgrade',
  templateUrl: './upgrade.component.html',
  styleUrls: []
})

export class UpgradeComponent implements OnInit 
{
  back: string = "";
  showCloseDownAccount: boolean = false;

  //
  // Construct.
  //
  constructor(private activatedRoute: ActivatedRoute) { }

  //
  // OnInit...
  //
  ngOnInit() {
    // subscribe to router event
    this.activatedRoute.queryParams.subscribe((params: Params) => {

      // Set the back
      if(params['back']) 
      {
        this.back = params['back'];
      }

    });

  }

  //
  // Cancel account
  //
  cancelAccount()
  {
    this.showCloseDownAccount = true;
  }

  //
  // Cancel account
  //
  doCancelAccount() 
  {
    this.showCloseDownAccount = false;
  }
}

/* End File */