//
// Date: 11/10/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Symbol } from './symbol';

export class WatchlistItems 
{
  public Id: string;
  public Symbol: Symbol;

  //
  // Construct
  //
  constructor(Id: string, Symb: Symbol) {
    this.Id = Id;
    this.Symbol = Symb;
  }  
}

/* End File */