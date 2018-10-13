//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { MeService } from '../../../providers/http/me.service';

declare var Stripe: any;

@Component({
  selector: 'app-settings-account-card',
  templateUrl: './card.component.html',
  styleUrls: []
})

export class CardComponent implements OnInit 
{
  // Passed vars
  @Input() showOverlay: boolean;
  @Output() onClose = new EventEmitter<boolean>();

  // Vars
  errMsg: string = "";
  couponMsg: string = "";
  btnDisabled: boolean = false;
  showCouponField: boolean = false;  
  form: CreditCardForm = new CreditCardForm();

  //
  // Construct.
  //
  constructor(private meService: MeService) {
    // Get cached data
    //this.userProfile = this.stateService.GetSettingsUserProfile();
  }

  // 
  // NG Init.
  //
  ngOnInit() 
  {
    this.form.Number = '4242424242424242';
    this.form.CVC = '123';
    this.form.ExpMonth = '12';
    this.form.ExpYear = '2019';
    this.form.ZipCode = '97132';               
  }

  //
  // Do show coupon.
  //
  doShowCoupon()
  {
    this.showCouponField = true;
  }

  //
  // Get credit card token.
  //
  getCreditCardToken() 
  {
    this.errMsg = "";
    this.btnDisabled = true;

    // Get Stripe token
    Stripe.card.createToken({
      number: this.form.Number,
      cvc: this.form.CVC,
      exp_month: this.form.ExpMonth,
      exp_year: this.form.ExpYear,
      address_zip: this.form.ZipCode
    }, 

    // Handle response.
    (status, response) => {

      // Set button
      this.btnDisabled = false;

      // Is this an error?
      if(status >= 300)
      {
        this.errMsg = response.error.message;
        return;
      }

      // Clear credit card data.
      this.form = new CreditCardForm();
      this.form.ExpMonth = '1';
      this.form.ExpYear = '2020';
    });
  }

  //
  // Cancel
  //
  doCancel()
  {
    this.showCouponField = false;

    // Clear credit card data.
    this.form = new CreditCardForm();
    this.form.ExpMonth = '1';
    this.form.ExpYear = '2020';

    // Close overlay
    this.showOverlay = false;
    this.onClose.emit(false);        
  }

}

class CreditCardForm 
{
  Number: string,
  CVC: string,
  ExpMonth: string,
  ExpYear: string,
  ZipCode: string,
  Coupon: string
}

/* End File */