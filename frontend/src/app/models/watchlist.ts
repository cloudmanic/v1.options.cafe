//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

export class Watchlist {
  
  constructor(
    public id: string,
    public name: string,    
    public items: string[]
  ){}
  
  //
  // Build build the data for emitting to the app. 
  //
  public static buildForEmit(data) : Watchlist {
    return new Watchlist(data.Id, data.Name, data.List);  
  }
  
}

/* End File */