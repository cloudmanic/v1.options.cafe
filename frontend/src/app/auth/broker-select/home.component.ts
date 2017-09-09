//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { NgForm }   from '@angular/forms';
import { Router } from '@angular/router';

declare var app_server: any;

@Component({
  selector: 'broker-select',
  templateUrl: './home.component.html'
})

export class AuthBrokerSelectComponent implements OnInit {

  broker: string = "tradier"

  constructor(private router: Router) { }

  //
  // On Init
  //
  ngOnInit() {
    
    // No user id redirect
    if(! localStorage.getItem('user_id').length)
    {
      this.router.navigate([ '/login' ]);
    }

    // No access token redirect
    if(! localStorage.getItem('access_token').length)
    {
      this.router.navigate([ '/login' ]);
    }
    
  }

  //
  // Login submit.
  //
  onSubmit(form: NgForm) {
    
    // Switch based on broker selected
    switch(form.value["field-broker"])
    {
      case 'tradier':
        window.location.href = app_server + '/tradier/authorize?user=' + localStorage.getItem('user_id');
      break;
    }
    
  }
  
}

/* End File */
