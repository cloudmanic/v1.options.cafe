import { WatchlistItems } from './watchlist-items';

export class Watchlist {
  constructor(
    public id: string,
    public name: string,    
    public items: WatchlistItems[]
  ){}
}
