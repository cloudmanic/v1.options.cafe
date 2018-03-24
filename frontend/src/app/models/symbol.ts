//
// Date: 11/7/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

//
// Symbol Model
//
export class Symbol 
{
  Id: number;
  Name: string;
  ShortName: string;
  Type: string;  
  OptionDetails: OptionDetails;

  //
  // Create a new object
  //
  New(Id: number, Name: string, ShortName: string, Type: string, OptionDetails: OptionDetails): Symbol 
  {
    let obj = new Symbol();

    obj.Id = Id;
    obj.Name = Name;
    obj.ShortName = ShortName;
    obj.Type = Type;
    obj.OptionDetails = OptionDetails;

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
    this.OptionDetails = new OptionDetails(json["option_details"].symbol, new Date(json["option_details"].expire), json["option_details"].strike, json["option_details"].type);
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

//
// Option details. 
//
export class OptionDetails {
  public Symbol: string;
  public Expire: Date;
  public Strike: number;
  public Type: string;

  //
  // Constructor
  //
  constructor(Symbol: string, Expire: Date, Strike: number, Type: string) {
    this.Symbol = Symbol;
    this.Expire = Expire;
    this.Strike = Strike;
    this.Type = Type;
  }  
}

/* End Find */