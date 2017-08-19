import { Component, OnInit } from '@angular/core';
import { NgForm }   from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router } from '@angular/router';

declare var app_server: any;

interface LoginResponse {
  status: number, 
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
  submitBtn = "Log In";

  constructor(private http: HttpClient, private router: Router) { }

  ngOnInit() {}
  
  //
  // Login submit.
  //
  onSubmit(form: NgForm) {

    // Clear post error.
    this.errorMsg = "";
    
    // Update submit button
    this.submitBtn = "Posting...";

    // Make the the HTTP request:
    this.http.post<LoginResponse>(app_server + '/login', form.value).subscribe(
      
      // Success
      data => {
        
        // Store access token in local storage. 
        localStorage.setItem('access_token', data.access_token); 
        
        // See if we have a broker or not.
        if(data.broker_count == 0)
        {
          this.router.navigate([ '/broker-select' ]);
        } else
        {
          this.router.navigate([ '/' ]);
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
          // The backend returned an unsuccessful response code.
          //console.log(`Backend returned code ${err.status}, body was: ${err.error}`);
          
          // Print error message
          this.errorMsg = err.error.error;
        }
        
      }
        
    );   
  
  }

}
