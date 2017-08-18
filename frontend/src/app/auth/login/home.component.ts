import { Component, OnInit } from '@angular/core';
import { NgForm }   from '@angular/forms';
import { HttpClient } from '@angular/common/http';

declare var app_server: any;

@Component({
  selector: 'app-auth-login',
  templateUrl: './home.component.html'
})

export class AuthLoginComponent implements OnInit {

  constructor(private http: HttpClient) { }

  ngOnInit() {
    
   
    
  }
  
  onSubmit(form: NgForm) {

    console.log(form.value);

    // Make the the HTTP request:
    this.http.post(app_server + '/login', form.value).subscribe(data => {

      console.log(data);

    });   
  
  }

}
