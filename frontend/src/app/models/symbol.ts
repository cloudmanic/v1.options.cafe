//
// Date: 11/7/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

export class Symbol {
  constructor(
    public Name: string,
    public Description: string  
  ){}

  //
  // Build object for emitting to the app.
  //
  public static buildForEmit(data) : Symbol[]  {

    let symbols = [];

    for(let i = 0; i < data.length; i++)
    {
      symbols.push(new Symbol(data[i].Name, data[i].Description));
    }    

    return symbols; 
  }  
}

/* End Find */