//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { NgForm }   from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router, ActivatedRoute, Params } from '@angular/router';
import { environment } from '../../../environments/environment';

interface LoginResponse {
  status: number, 
  user_id: number,
  error: string,
  access_token: string,
  broker_count: number
}

@Component({
  selector: 'app-auth-login',
  templateUrl: './home.component.html'
})

export class AuthLoginComponent implements OnInit {

  errorMsg = "";
  successMsg = "";
  submitBtn = "Log In";
  returnUrl: "/";

  constructor(private http: HttpClient, private router: Router, private activatedRoute: ActivatedRoute) { }

  //
  // OnInit...
  //
  ngOnInit() {
    
    // Remove local storage
    localStorage.removeItem('user_id');    
    localStorage.removeItem('access_token');
    localStorage.removeItem('active_account');    
    
    // subscribe to router event
    this.activatedRoute.queryParams.subscribe((params: Params) => {
      this.successMsg = params['success'];
    });
    
    // get return url from route parameters or default to '/'
    this.returnUrl = this.activatedRoute.snapshot.queryParams['returnUrl'] || '/';    
    
  }
  
  //
  // Login submit.
  //
  onSubmit(form: NgForm) {

    // Clear post error.
    this.errorMsg = "";
    
    // Update submit button
    this.submitBtn = "Posting...";

    // Make the the HTTP request:
    this.http.post<LoginResponse>(environment.app_server + '/login', form.value).subscribe(
      
      // Success
      data => {
        
        // Store access token in local storage. 
        localStorage.setItem('user_id', data.user_id.toString());
        localStorage.setItem('access_token', data.access_token); 
        
        // See if we have a broker or not.
        if(data.broker_count == 0)
        {
          this.router.navigate([ '/broker-select' ]);
        } else
        {
          this.router.navigate([ this.returnUrl ]);
        }

      },
      
      // Error
      (err: HttpErrorResponse) => {

        // Change button back.
        this.submitBtn = "Log In";

        if (err.error instanceof Error) 
        {
          // A client-side or network error occurred. Handle it accordingly.
          console.log('An error occurred:', err.error.message);
        } else 
        { 
          // Print error message
          var json = JSON.parse(err.error); // Bug....Angular 4.4.4
          this.errorMsg = json.error;
        }
        
      }
        
    );   
  
  }

}
