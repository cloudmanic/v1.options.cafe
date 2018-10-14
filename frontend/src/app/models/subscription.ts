//
// Date: 10/11/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';

export class Subscription 
{
  Name: string;
  Amount: number;
  Status: string;
  TrialDays: number;
  TrialStart: Date;
  TrialEnd: Date;
  CurrentPeriodStart: Date;
  CurrentPeriodEnd: Date;
  BillingInterval: string;
  CardLast4: string;
  CardBrand: string;
  CardExpireMonth: number;
  CardExpireYear: number;
  CouponName: string;
  CouponCode: string;
  CouponAmountOff: number;
  CouponPercentOff: number;
  CouponDuration: string;  

  //
  // Build from JSON.
  //
  fromJson(json: Object): Subscription 
  {
    let result = new Subscription();

    if (!json) {
      return result;
    }

    // Set data.
    result.Name = json["name"];
    result.Amount = json["amount"];
    result.Status = json["status"];
    result.TrialDays = json["trial_days"];
    result.TrialStart = moment(json["trial_start"]).toDate();
    result.TrialEnd = moment(json["trial_end"]).toDate();
    result.CurrentPeriodStart = moment(json["current_period_start"]).toDate();
    result.CurrentPeriodEnd = moment(json["current_period_end"]).toDate();
    result.BillingInterval = json["billing_interval"];
    result.CardLast4 = json["card_last_4"];
    result.CardBrand = json["card_brand"];
    result.CardExpireMonth = json["card_exp_month"];
    result.CardExpireYear = json["card_exp_year"];
    result.CouponName = json["coupon_name"];
    result.CouponCode = json["coupon_code"];
    result.CouponAmountOff = json["coupon_amount_off"];
    result.CouponPercentOff = json["coupon_percent_off"];
    result.CouponDuration = json["coupon_duration"];     

    // Return happy
    return result;
  }
}



/* End File */
