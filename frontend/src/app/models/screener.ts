//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { ScreenerResult } from './screener-result';

//
// Screener
//
export class Screener 
{
  Id: number;
  Name: string;
  Strategy: string;
  Symbol: string;
  Items: ScreenerItem[];
  Results: ScreenerResult[];

  //
  // Build from JSON list.
  //
  fromJsonList(json: Object[]): Screener[] {
    let result = [];

    if (!json) 
    {
      return result;
    }

    for (let i = 0; i < json.length; i++) 
    {
      result.push(new Screener().fromJson(json[i]));
    }

    // Return happy
    return result;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): Screener 
  {
    let obj = new Screener();

    obj.Id = json["id"];
    obj.Name = json["name"];
    obj.Strategy = json["strategy"];
    obj.Symbol = json["symbol"];
    obj.Items = [];

    // Add in the legs.
    for (let i = 0; i < json["items"].length; i++)
    {
      obj.Items.push(new ScreenerItem().fromJson(json["items"][i]));
    }

    return obj;
  }  
}

//
// Screener Item
//
export class ScreenerItem 
{
  Id: number;
  ScreenerId: number;
  Key: string;
  Operator: string;
  ValueString: string;
  ValueNumber: number;

  //
  // Json to Object.
  //
  fromJson(json: Object): ScreenerItem {
    let obj = new ScreenerItem();

    obj.Id = json["id"];
    obj.Key = json["key"];
    obj.ScreenerId = json["screener_id"];
    obj.Operator = json["operator"];
    obj.ValueString = json["value_string"];
    obj.ValueNumber = json["value_number"];

    return obj;
  } 
}

/* End File */