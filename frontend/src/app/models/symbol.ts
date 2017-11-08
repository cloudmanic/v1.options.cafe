//
// Date: 11/7/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

export class Symbol {
  constructor(
    public ShortName: string,
    public Name: string  
  ){}

  //
  // Build object for emitting to the app.
  //
  public static buildForEmit(data) : Symbol[]  {

    let symbols = [];

    for(let i = 0; i < data.length; i++)
    {
      symbols.push(new Symbol(data[i].ShortName, data[i].Name));
    }    

    return symbols; 
  }  
}

/* End Find */