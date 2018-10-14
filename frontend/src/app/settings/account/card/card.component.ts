//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { Coupon } from '../../../models/coupon';
import { MeService } from '../../../providers/http/me.service';
import { StateService } from '../../../providers/state/state.service';

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
  coupon: Coupon = new Coupon(); 
  form: CreditCardForm = new CreditCardForm();

  //
  // Construct.
  //
  constructor(private stateService: StateService, private meService: MeService) { }

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
    this.form.Coupon = 'Go7s5ivW';               
  }

  //
  // Do show coupon.
  //
  doShowCoupon()
  {
    this.showCouponField = true;
  }

  //
  // Verify coupon
  //
  verifyCoupon()
  {
    this.errMsg = "";

    // Ajax call to verify coupon.
    this.meService.getVerifyCoupon(this.form.Coupon).subscribe((res) => {
      this.coupon = res;

      // If 100% off no need for a credit card.
      if((this.coupon.PercentOff >= 100) && (this.coupon.Valid))
      {
        this.meService.applyCoupon(this.form.Coupon).subscribe((res) => {
          this.doCancel();
        });
      } else 
      {
        this.couponMsg = this.coupon.PercentOff + "% off your subscription.";
      }
    },

    // Error
    (err: HttpErrorResponse) => {
      this.errMsg = "Invalid coupon code.";
    });
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

      // Send the token to the server
      this.meService.updateCreditCard(response.id, this.form.Coupon).subscribe((res) => {

        // Close overlay
        this.showOverlay = false;
        this.onClose.emit(false); 

        // Show success notice
        this.stateService.SiteSuccess.emit("Your credit card has been successfully updated.");

        // Clear credit card data.
        this.form = new CreditCardForm();
        this.form.ExpMonth = '1';
        this.form.ExpYear = '2020';  

      });

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