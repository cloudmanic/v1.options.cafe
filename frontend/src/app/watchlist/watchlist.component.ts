import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { QuoteService } from '../services/quote.service';
import { BrokerService } from '../services/broker.service';
import { MarketQuote } from '../contracts/market-quote';
import { Watchlist } from '../contracts/watchlist';
import { WatchlistItems } from '../contracts/watchlist-items';

@Component({
  selector: 'oc-watchlist',
  templateUrl: './watchlist.component.html'
})
export class WatchlistComponent implements OnInit {
  quotes = {};
  watchlist = new Watchlist('', '', []);

  //
  // Constructor....
  //
  constructor(private quotesService: QuoteService, private broker: BrokerService, private changeDetect: ChangeDetectorRef) { }

  //
  // OnInit....
  //
  ngOnInit() {

    // Subscribe to data updates from the broker - Watchlist
    this.broker.watchlistPushData.subscribe(data => {
      this.watchlist = data;
      this.changeDetect.detectChanges();
    });

    // Subscribe to data updates from the quotes - Market Quotes
    this.quotesService.marketQuotePushData.subscribe(data => {
      
      this.quotes[data.symbol] = data;
      this.changeDetect.detectChanges();
      
    });

  }

}
