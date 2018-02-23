//
// Date: 2/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { StateService } from '../../../providers/state/state.service';

@Component({
  selector: '[app-trading-ivr]',
  templateUrl: './ivr.component.html'
})

export class IvrComponent implements OnInit {
  
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