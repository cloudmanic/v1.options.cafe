import { Component, OnInit } from '@angular/core';
import { NgForm }   from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router } from '@angular/router';

declare var app_server: any;

@Component({
  selector: 'app-auth-forgot-password',
  templateUrl: './home.component.html'
})
export class AuthForgotPasswordComponent implements OnInit {

  errorMsg = "";
  submitBtn = "Reset Password";

  constructor(private http: HttpClient, private router: Router) { }

  ngOnInit() {}
  
  //
  // Reset submit.
  //
  onSubmit(form: NgForm) {
        
    // Clear post error.
    this.errorMsg = "";
    
    // Update submit button
    this.submitBtn = "Posting...";

    // Make the the HTTP request:
    this.http.post(app_server + '/forgot-password', form.value).subscribe(
      
      // Success redirect to login
      data => {
        this.router.navigate([ '/login' ], { queryParams: { success: "Please check your email for next steps." } });
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