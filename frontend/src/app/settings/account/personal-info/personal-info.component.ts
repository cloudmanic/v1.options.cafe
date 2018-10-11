//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Me } from '../../../models/me';
import { Component, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { MeService } from '../../../providers/http/me.service';

@Component({
  selector: '[app-settings-account-personal-info]',
  templateUrl: './personal-info.component.html',
  styleUrls: []
})

export class PersonalInfoComponent implements OnInit {

  tmpUserProfile: Me = new Me();
  userProfile: Me = new Me();
  showEditProfile: boolean = false;
  showEditPassword: boolean = false;

  // Validation 
  firstNameError = ""
  lastNameError = ""
  emailError = ""    

  //
  // Construct.
  //
  constructor(private meService: MeService) { }

  // 
  // NG Init.
  //
  ngOnInit() 
  { 
    // Load page data.
    this.getUserProfile();
  }

  //
  // Get user data.
  //
  getUserProfile()
  {
    // Ajax call to get user data.
    this.meService.getProfile().subscribe((res) => {
      this.userProfile = res;
    });
  }

  //
  // Edit profile.
  //
  doShowEditProfile() 
  {
    // Create temp object
    this.tmpUserProfile = new Me().setFromObj(this.userProfile);
    this.showEditProfile = true;
  }

  //
  // Save Edit profile.
  //
  doSaveEditProfile() 
  {
    let yesError = false;

    // Clear Validation 
    this.firstNameError = "";
    this.lastNameError = "";
    this.emailError = ""; 

    // Validate - First Name
    if (this.userProfile.FirstName.length <= 0)
    {
      yesError = true;
      this.firstNameError = "First name field is required.";
    }

    // Validate - Last Name
    if (this.userProfile.LastName.length <= 0) 
    {
      yesError = true;      
      this.lastNameError = "Last name field is required.";
    }

    // Validate - Email
    if (this.userProfile.Email.length <= 0) 
    {
      yesError = true;
      this.emailError = "Email field is required.";
    }

    if(yesError)
    {
      return;
    }

    // Ajax call to save the profile data.
    this.meService.saveProfile(this.userProfile).subscribe(
      
      // Success
      (res) => {
        this.userProfile = res;
        this.showEditProfile = false;
      }, 

      // Error
      (err: HttpErrorResponse) => {
        this.emailError = err.error.errors.email;
      }

    );
  }

  //
  // Cancel profile.
  //
  doCancelEditProfile() 
  {
    // Clear Validation 
    this.firstNameError = "";
    this.lastNameError = "";
    this.emailError = ""; 

    // Reset values    
    this.userProfile = new Me().setFromObj(this.tmpUserProfile);
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