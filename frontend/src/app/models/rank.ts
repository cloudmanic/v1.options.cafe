//
// Date: 7/15/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

export class Rank 
{
  public Rank30: number;
  public Rank60: number;  
  public Rank90: number;
  public Rank365: number;

  //
  // Constructor...
  //
  constructor(
    Rank30: number,
    Rank60: number,
    Rank90: number,
    Rank365: number
  ){
    this.Rank30 = Rank30;
    this.Rank60 = Rank60;
    this.Rank90 = Rank90;    
    this.Rank365 = Rank365;
  }  
}

/* End File */