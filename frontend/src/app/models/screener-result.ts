//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Symbol } from './symbol';

//
// OptionsChain
//
export class ScreenerResult 
{
  Debit: number;
  Credit: number;  
  MidPoint: number;
  PutPrecentAway: number;
  CallPrecentAway: number;  
  Legs: Symbol[];

  //
  // Build from JSON list.
  //
  fromJsonList(json: Object[]): ScreenerResult[] 
  {
    let result = [];

    if (!json) {
      return result;
    }

    for (let i = 0; i < json.length; i++) 
    {
      result.push(new ScreenerResult().fromJson(json[i]));
    }

    // Return happy
    return result;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): ScreenerResult 
  {
    let obj = new ScreenerResult();

    obj.Debit = json["debit"];
    obj.Credit = json["credit"];
    obj.MidPoint = json["midpoint"];
    obj.PutPrecentAway = json["put_percent_away"];
    obj.CallPrecentAway = json["call_percent_away"];    
    obj.Legs = [];

    // Add in the legs.
    for (let i = 0; i < json["legs"].length; i++)
    {
      obj.Legs.push(new Symbol().fromJson(json["legs"][i]));
    }

    return obj;
  }  
}

/* End File */