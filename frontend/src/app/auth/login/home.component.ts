//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { NgForm }  from '@angular/forms';
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

interface GoogleSessionResponse {
  session_secret: string
}

@Component({
  selector: 'app-auth-login',
  templateUrl: './home.component.html'
})

export class AuthLoginComponent implements OnInit {

  errorMsg = "";
  successMsg = "";
  submitBtn = "Login";
  returnUrl: "/";
  googleLoginState: boolean = false;

  //
  // Construct.
  //
  constructor(private http: HttpClient, private router: Router, private activatedRoute: ActivatedRoute) { }

  //
  // OnInit...
  //
  ngOnInit() {

    // Remove access token on server
    let access_token = localStorage.getItem('access_token');

    if(access_token && access_token.length)
    {
      this.http.get(environment.app_server + '/oauth/logout?access_token=' + access_token).subscribe();
    }

    // Remove local storage
    localStorage.removeItem('user_id');
    localStorage.removeItem('redirect');
    localStorage.removeItem('broker_new_id');    
    localStorage.removeItem('access_token');
    localStorage.removeItem('active_account');
    localStorage.removeItem('active_watchlist');    
    
    // subscribe to router event
    this.activatedRoute.queryParams.subscribe((params: Params) => {

      // Success
      this.successMsg = params['success'];

      // See if we have a google_auth_success for a successful google login.
      if(params['google_auth_success']) 
      {
        this.processPostGoogleLogin();
      }

      // See if our google login failed to login.
      if(params['google_auth_failed']) 
      {
        this.errorMsg = "Looks like you haven’t linked up that Google account yet.";
      }
    });
    
    // get return url from route parameters or default to '/'
    this.returnUrl = this.activatedRoute.snapshot.queryParams['returnUrl'] || '/';    
  }

  //
  // Process post Google login.
  //
  processPostGoogleLogin()
  {
    let sessionKey = localStorage.getItem('google_auth_session_key');
    let sessionSecret = localStorage.getItem('google_auth_session_secret');

    // Remove keys we do not need.
    localStorage.removeItem('google_auth_session_key');
    localStorage.removeItem('google_auth_session_secret');


    let post = { 
       session_key: sessionKey, 
       session_secret: sessionSecret, 
       grant_type: "password", 
       client_id: environment.client_id
     }

    // Make the the HTTP request:
    this.http.post<LoginResponse>(environment.app_server + '/oauth/google/token', post).subscribe(

      // Success
      data => {
        // Store access token in local storage. 
        localStorage.setItem('user_id', data.user_id.toString());
        localStorage.setItem('access_token', data.access_token);

        // See if we have a broker or not.
        if (data.broker_count == 0) {
          this.router.navigate(['/broker-select']);
        } else {
          this.router.navigate([this.returnUrl]);
        }
      },

      // Error
      (err: HttpErrorResponse) => {
        this.errorMsg = "Unable to login via your Google account.";
        console.log('An error occurred:', err);
      }

    );
  }

  //
  // Do Google login.
  //
  googleLogin()
  {
    this.googleLoginState = true;
    let sessionKey = this.getRandomString();
    localStorage.setItem('google_auth_session_key', sessionKey);    

    // Make the the HTTP request:
    this.http.post<GoogleSessionResponse>(environment.app_server + '/oauth/google/session', { session_key: sessionKey, type: 'login', redirect: environment.site_url + '/login' }).subscribe(

      // Success
      data => {
        // Store session secret. 
        localStorage.setItem('google_auth_session_secret', data.session_secret.toString());

        // Redirect to start the Google login
        window.location.href = environment.app_server + '/oauth/google?session_key=' + sessionKey;
      },

      // Error
      (err: HttpErrorResponse) => {
        console.log('An error occurred:', err);
      }

    );
  }
  
  //
  // Login submit.
  //
  onSubmit(form: NgForm) {

    if(this.googleLoginState)
    {
      return;
    }

    // Clear post error.
    this.errorMsg = "";
    
    // Update submit button
    this.submitBtn = "Logging In...";

    // Add oauth stuff
    form.value.grant_type = "password";
    form.value.client_id = environment.client_id;

    // Make the the HTTP request:
    this.http.post<LoginResponse>(environment.app_server + '/oauth/token', form.value).subscribe(
      
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
        this.submitBtn = "Login";

        if (err.error instanceof Error) 
        {
          // A client-side or network error occurred. Handle it accordingly.
          console.log('An error occurred:', err.error);
        } else 
        { 
          // Print error message
          this.errorMsg = err.error.error;
        }
        
      }
        
    );   
  
  }

  //
  // Make a random string to manage the google access token request.
  //
  getRandomString()
  {
    let text = "";
    let possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (var i = 0; i < 20; i++)
    {
      text += possible.charAt(Math.floor(Math.random() * possible.length));   
    }

    return text;
  }

}
