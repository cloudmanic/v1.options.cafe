//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-expired',
  templateUrl: './expired.component.html',
  styleUrls: []
})

export class ExpiredComponent implements OnInit 
{
  showCloseDownAccount: boolean = false;

  //
  // Construct.
  //
  constructor() { }

  //
  // OnInit...
  //
  ngOnInit() { }

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