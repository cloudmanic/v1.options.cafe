//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Symbol } from './symbol';
import { WatchlistItems } from './watchlist-items';

export class Watchlist 
{
  public Id: number;
  public Name: string;    
  public Symbols: WatchlistItems[];

  //
  // Construct
  //
  constructor(Id: number, Name: string, Items: WatchlistItems[]) {
    this.Id = Id;
    this.Name = Name;
    this.Symbols = Items;
  } 
 
  //
  // Build build the data for emitting to the app. 
  //
  public static buildForEmit(data) : Watchlist {

    let symbs = [];

    // Build Items
    for(let i = 0; i < data.symbols.length; i++)
    {
      symbs.push(new WatchlistItems(data.symbols[i].id, new Symbol().fromJson(data.symbols[i].symbol)));
    }

    // Return happy.
    return new Watchlist(data.id, data.name, symbs);
  }
  
}

/* End File */