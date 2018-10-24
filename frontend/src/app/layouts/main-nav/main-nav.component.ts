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
  showCentcom: boolean = false;
  routeData: any;  
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
    // Route data
    this.routeData = this.route.data.subscribe(v => {
      this.section = v.section;
    });

    // See if we should show the centcom link
    if(localStorage.getItem('user_id_centcom') && (localStorage.getItem('user_id_centcom').length > 0)) 
    {
      this.showCentcom = true;
    }
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

  //
  // Log back into centcom
  //
  logBackIntoCentcom()
  {
    // Remove local storage
    localStorage.removeItem('user_id');
    localStorage.removeItem('redirect');
    localStorage.removeItem('broker_new_id');
    localStorage.removeItem('access_token');
    localStorage.removeItem('active_account');
    localStorage.removeItem('active_watchlist');

    // Store access token in local storage. 
    localStorage.setItem('user_id', localStorage.getItem('user_id_centcom'));
    localStorage.setItem('access_token', localStorage.getItem('access_token_centcom'));
    localStorage.setItem('active_account', localStorage.getItem('active_account_centcom'));
    localStorage.setItem('active_watchlist', localStorage.getItem('active_watchlist_centcom'));

    // Remove local storage
    localStorage.removeItem('user_id_centcom');
    localStorage.removeItem('access_token_centcom');
    localStorage.removeItem('active_account_centcom');
    localStorage.removeItem('active_watchlist_centcom');

    // Redirect to Centcom
    window.location.href = '/centcom/users';
  }

}

/* End File */