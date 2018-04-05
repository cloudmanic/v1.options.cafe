//
// Date: 4/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Symbol } from './symbol';

//
// OptionsChain
//
export class OptionsChain 
{
  Underlying: string;
  ExpirationDate: Date;
  Calls: OptionsChainItem[];
  Puts: OptionsChainItem[];


  //
  // Json to Object.
  //
  fromJson(json: Object): OptionsChain 
  {
    let obj = new OptionsChain();

    obj.Underlying = json["underlying"];
    obj.ExpirationDate = moment(json["expiration_date"]).toDate();
    obj.Calls = [];
    obj.Puts = [];

    // Add in the calls.
    for (let i = 0; i < json["calls"].length; i++)
    {
      obj.Calls.push(new OptionsChainItem().fromJson(json["calls"][i]));
    }

    // Add in the puts.
    for (let i = 0; i < json["puts"].length; i++) {
      obj.Puts.push(new OptionsChainItem().fromJson(json["puts"][i]));
    }   

    return obj;
  }  
}

//
// OptionsChainItem
//
export class OptionsChainItem
{
  Underlying: string;
  Symbol: string;
  OptionType: string;
  Description: number;
  Strike: number;
  ExpirationDate: Date;
  Last: number;
  Change: number;
  ChangePercentage: number;
  Volume: number;
  AverageVolume: number;
  LastVolume: number;
  Open: number;
  High: number;
  Low: number;
  Close: number;
  Bid: number;
  BidSize: number;
  Ask: number;
  AskSize: number;
  OpenInterest: number;

  //
  // Json to Object.
  //
  fromJson(json: Object): OptionsChainItem 
  {
    let obj = new OptionsChainItem();
  
    obj.Underlying = json["underlying"];
    obj.Symbol = json["symbol"];
    obj.OptionType = json["option_type"];
    obj.Description = json["description"];
    obj.Strike = json["strike"];
    obj.ExpirationDate = moment(json["expiration_date"]).toDate();
    obj.Last = json["last"];
    obj.Change = json["change"];
    obj.ChangePercentage = json["change_percentage"];
    obj.Volume = json["volume"];
    obj.AverageVolume = json["average_volume"];
    obj.LastVolume = json["last_volume"];
    obj.Open = json["open"];
    obj.High = json["high"];
    obj.Low = json["low"];
    obj.Close = json["close"];
    obj.Bid = json["bid"];
    obj.BidSize = json["bid_size"];
    obj.Ask = json["ask"];
    obj.AskSize = json["ask_size"];
    obj.OpenInterest = json["open_interest"];

    return obj;
  }
}

/* End File */