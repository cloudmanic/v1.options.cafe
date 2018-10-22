//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Observable';
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: []
})

export class UsersComponent implements OnInit 
{
  users: Object[] = [];

  //
  // Constructor
  //
  constructor(private http: HttpClient) { }

  //
  // NgInit
  //
  ngOnInit() 
  { 
    this.getUsers();
  }

  //
  // Get users 
  //
  getUsers() 
  {
    this.http.get<Object[]>(environment.app_server + '/api/admin/users').subscribe((data) => {
      this.users = data;
    });
  }

  //
  // Login as user.
  //
  loginAsUser(user: Object) 
  {
    alert('adsf');
  }

}

/* End File */