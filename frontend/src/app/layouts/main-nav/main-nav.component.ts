//
// Date: 9/7/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//


declare var groove: any;

import { ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-main-nav',
  templateUrl: './main-nav.component.html'
})

export class MainNavComponent implements OnInit 
{
  section: string = 'trading';

  //
  // Construct.
  //
  constructor(private route: ActivatedRoute) { }

  //
  // OnInit
  //
  ngOnInit() 
  {
    this.routeData = this.route.data.subscribe(v => {
      this.section = v.section;
    });
  }

  //
  // OnDestory
  //
  ngOnDestroy() 
  {
    this.routeData.unsubscribe();
  }
  //
  // Clicked on help.
  //
  onHelpClick() 
  { 
    groove.widget('open');
  }

}

/* End File */