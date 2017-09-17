//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

export class MarketStatus {
  
  //
  // Construct.
  //
  constructor(
    public description: string,
    public state: string
  ){}
  
  //
  // Build object for emitting to the app.
  //
  public static buildForEmit(data) : MarketStatus  {
    return new MarketStatus(data.State, data.Description); 
  }
}

/* End Files */