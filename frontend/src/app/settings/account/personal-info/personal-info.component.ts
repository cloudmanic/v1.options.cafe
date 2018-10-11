//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Me } from '../../../models/me';
import { HttpErrors } from '../../../models/http-errors';
import { Component, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { MeService } from '../../../providers/http/me.service';
import { StateService } from '../../../providers/state/state.service';

@Component({
  selector: '[app-settings-account-personal-info]',
  templateUrl: './personal-info.component.html',
  styleUrls: []
})

export class PersonalInfoComponent implements OnInit 
{
  httpErrors: HttpErrors = new HttpErrors();
  tmpUserProfile: Me = new Me();
  userProfile: Me = new Me();
  showEditProfile: boolean = false;
  showEditPassword: boolean = false;

  // Password 
  currentPass = "";
  newPass = "";
  newPassConfirm = "";
  passError = "";    

  //
  // Construct.
  //
  constructor(private stateService: StateService, private meService: MeService) 
  { 
    // Get cached data
    this.userProfile = this.stateService.GetSettingsUserProfile();
  }

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
      this.stateService.SetSettingsUserProfile(res);
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
    this.httpErrors = new HttpErrors();

    // Validate - First Name
    if (this.userProfile.FirstName.length <= 0)
    {
      yesError = true;
      this.httpErrors.FirstName = "First name field is required.";
    }

    // Validate - Last Name
    if (this.userProfile.LastName.length <= 0) 
    {
      yesError = true;      
      this.httpErrors.LastName = "Last name field is required.";
    }

    // Validate - Email
    if (this.userProfile.Email.length <= 0) 
    {
      yesError = true;
      this.httpErrors.Email = "Email field is required.";
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
        this.stateService.SetSettingsUserProfile(res);

        // Show success notice
        this.stateService.SiteSuccess.emit("Your profile has been successfully updated.");        
      }, 

      // Error
      (err: HttpErrorResponse) => {
        this.httpErrors = new HttpErrors().fromJson(err.error.errors);
      }

    );
  }

  //
  // Cancel profile.
  //
  doCancelEditProfile() 
  {
    // Clear Validation 
    this.httpErrors = new HttpErrors();

    // Reset values    
    this.userProfile = new Me().setFromObj(this.tmpUserProfile);
    this.showEditProfile = false;
  }

  //
  // Edit Password.
  //
  doShowEditPassword() {
    this.passError = "";
    this.currentPass = "";
    this.newPass = "";
    this.newPassConfirm = "";     
    this.showEditPassword = true;
  }

  //
  // Save Edit profile.
  //
  doSaveEditPassword() {

    this.passError = "";

    // Confirm passwords
    if(this.newPass != this.newPassConfirm)
    {
      this.passError = "New passwords did not match.";
      return;
    }

    // Ajax call to save the password.
    this.meService.restPassword(this.currentPass, this.newPass).subscribe(

      // Success
      (res) => {
        this.showEditPassword = false;

        // Show success notice
        this.stateService.SiteSuccess.emit("Your password has been successfully updated.");
      },

      // Error
      (err: HttpErrorResponse) => {
        this.passError = err.error.error;
      }

    );

    //this.showEditPassword = false;
  }

  //
  // Cancel profile.
  //
  doCancelEditPassword() { 
    this.passError = "";  
    this.currentPass = "";
    this.newPass = "";
    this.newPassConfirm = "";     
    this.showEditPassword = false;
  }



}

/* End File */