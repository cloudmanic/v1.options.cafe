//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Symbol } from './symbol';
import { WatchlistItems } from './watchlist-items';

export class Watchlist 
{
  Id: number;
  Name: string;    
  Symbols: WatchlistItems[];

  //
  // Json to Object.
  //
  fromJson(json: Object): Watchlist {
    let wl = new Watchlist();
    wl.Id = json["id"];
    wl.Name = json["name"];
    wl.Symbols = [];

    // Build Items
    if (json["symbols"]) 
    {
      for (let i = 0; i < json["symbols"].length; i++) {
        wl.Symbols.push(new WatchlistItems(json["symbols"][i].id, new Symbol().fromJson(json["symbols"][i].symbol)));
      }
    }    

    return wl;
  }  
}

/* End File */