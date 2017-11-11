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

  //
  // Constructor
  //
  constructor(Id: number, Name: string, ShortName: string) {
    this.Id = Id;
    this.Name = Name;
    this.ShortName = ShortName;
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
      symbols.push(new Symbol(data[i].id, data[i].name, data[i].short_name));
    }    

    return symbols; 
  }
}

/* End Find */