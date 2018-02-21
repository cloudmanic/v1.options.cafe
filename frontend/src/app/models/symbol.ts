//
// Date: 11/7/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

//
// Symbol Model
//
export class Symbol {
  public Id: number;
  public Name: string;
  public ShortName: string;
  public Type: string;  
  public OptionDetails: OptionDetails;

  //
  // Constructor
  //
  constructor(Id: number, Name: string, ShortName: string, Type: string, OptionDetails: OptionDetails) {
    this.Id = Id;
    this.Name = Name;
    this.ShortName = ShortName;
    this.Type = Type;    
    this.OptionDetails = OptionDetails;
  }

  //
  // Build object for emitting to the app.
  //
  public static buildForEmit(data) : Symbol[]  {

    let symbols = [];

    if(! data)
    {
      return symbols;      
    }

    for(let i = 0; i < data.length; i++)
    {
      symbols.push(new Symbol(
        data[i].id, 
        data[i].name,         
        data[i].short_name,
        data[i].type,         
        new OptionDetails(new Date(data[i].option_details.expire), data[i].option_details.strike, data[i].option_details.type)
       ));
    }    

    return symbols; 
  }
}

export class OptionDetails {
  public Expire: Date;
  public Strike: number;
  public Type: string;

  //
  // Constructor
  //
  constructor(Expire: Date, Strike: number, Type: string) {
    this.Expire = Expire;
    this.Strike = Strike;
    this.Type = Type;
  }  
}

/* End Find */