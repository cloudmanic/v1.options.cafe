//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';

@Component({
  selector: '[app-settings-account-personal-info]',
  templateUrl: './personal-info.component.html',
  styleUrls: []
})

export class PersonalInfoComponent implements OnInit {

  showEditProfile: boolean = false;
  showEditPassword: boolean = false;

  //
  // Constructor
  //
  constructor() { }

  // 
  // NG Init.
  //
  ngOnInit() { }

  //
  // Edit profile.
  //
  doShowEditProfile() {
    this.showEditProfile = true;
  }

  //
  // Save Edit profile.
  //
  doSaveEditProfile() {
    this.showEditProfile = false;
  }

  //
  // Cancel profile.
  //
  doCancelEditProfile() {
    this.showEditProfile = false;
  }

  //
  // Edit Password.
  //
  doShowEditPassword() {
    this.showEditPassword = true;
  }

  //
  // Save Edit profile.
  //
  doSaveEditPassword() {
    this.showEditPassword = false;
  }

  //
  // Cancel profile.
  //
  doCancelEditPassword() {
    this.showEditPassword = false;
  }



}

/* End File */