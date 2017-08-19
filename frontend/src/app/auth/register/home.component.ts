import { Component, OnInit } from '@angular/core';
import { NgForm }   from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router } from '@angular/router';

declare var app_server: any;

interface RegisterResponse {
  status: number, 
  error: string,
  access_token: string
}

@Component({
  selector: 'app-auth-register',
  templateUrl: './home.component.html'
})
export class AuthRegisterComponent implements OnInit {

  errorMsg = "";
  submitBtn = "Create Account";

  constructor(private http: HttpClient, private router: Router) { }

  ngOnInit() {}
  
  //
  // Register submit.
  //
  onSubmit(form: NgForm) {

    // Clear post error.
    this.errorMsg = "";
    
    // Update submit button
    this.submitBtn = "Saving...";

    // Make the the HTTP request:
    this.http.post<RegisterResponse>(app_server + '/register', form.value).subscribe(
      
      // Success
      data => {
        
        console.log(data);
        
        // Store access token in local storage. 
        localStorage.setItem('access_token', data.access_token); 
        
        // Redirect to broker select
        this.router.navigate([ '/broker-select' ]);

      },
      
      // Error
      (err: HttpErrorResponse) => {
        
        // Change button back.
        this.submitBtn = "Create Account";

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
