//
// Date: 11/7/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';

//
// Symbol Model
//
export class Symbol 
{
  Id: number;
  Name: string;
  ShortName: string;
  Type: string;  
  OptionUnderlying: string;
  OptionType: string;
  OptionExpire: Date;
  OptionStrike: number;

  //
  // Create a new object
  //
  New(Id: number, Name: string, ShortName: string, Type: string, underlying: string, optionType: string, expire: Date, strike: number): Symbol 
  {
    let obj = new Symbol();

    obj.Id = Id;
    obj.Name = Name;
    obj.ShortName = ShortName;
    obj.Type = Type;
    obj.OptionUnderlying = underlying;
    obj.OptionType = optionType;
    obj.OptionExpire = expire;
    obj.OptionStrike = strike;

    return obj;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): Symbol 
  {
    this.Id = json["id"];
    this.Name = json["name"];
    this.ShortName = json["short_name"];
    this.Type = json["type"];
    this.OptionUnderlying = json["option_underlying"];
    this.OptionType = json["option_type"];
    this.OptionExpire = moment(json["option_expire"]).toDate();
    this.OptionStrike = json["option_strike"];
    return this;
  }

  //
  // Build from JSON list.
  //
  fromJsonList(json: Object[]): Symbol[] {
    let result = [];

    if (!json) {
      return result;
    }

    for (let i = 0; i < json.length; i++) {
      result.push(new Symbol().fromJson(json[i]));
    }

    // Return happy
    return result;
  }
}

/* End Find */