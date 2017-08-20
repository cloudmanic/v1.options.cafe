import { Component, OnInit } from '@angular/core';
import { NgForm }   from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router, ActivatedRoute, Params } from '@angular/router';

declare var app_server: any;

@Component({
  selector: 'app-auth-reset-password',
  templateUrl: './home.component.html'
})

export class AuthResetPasswordComponent implements OnInit {

  hash = "";
  errorMsg = "";
  submitBtn = "Reset Password";

  constructor(private http: HttpClient, private router: Router, private activatedRoute: ActivatedRoute) { }

  //
  // OnInit...
  //
  ngOnInit() {
    
    // subscribe to router event
    this.activatedRoute.queryParams.subscribe((params: Params) => {
      this.hash = params['hash'];
    });
    
  }

  //
  // Reset submit.
  //
  onSubmit(form: NgForm) {
    
    // First make sure the passwords match.
    if(form.value.password != form.value.password_again)
    {
      this.errorMsg = "Opps, the passwords do not match.";
      return;
    }
    
    // Clear post error.
    this.errorMsg = "";
    
    // Update submit button
    this.submitBtn = "Saving...";

    // Add the hash to the post.
    form.value.hash = this.hash;

    // Make the the HTTP request:
    this.http.post(app_server + '/reset-password', form.value).subscribe(
      
      // Success redirect to login
      data => {
        this.router.navigate([ '/login' ], { queryParams: { success: "Your password was successfully reset." } });
      },
      
      // Error
      (err: HttpErrorResponse) => {

        // Change button back.
        this.submitBtn = "Reset Password";

        if(err.error instanceof Error) 
        {
          console.log('A client-side error occurred:', err.error.message);
        } else 
        {
          this.errorMsg = err.error.error;
        }
        
      }
        
    );   
  
  }

}

/* End File */