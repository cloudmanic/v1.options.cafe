//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Rx';
import { Subject } from 'rxjs/Subject';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Screener } from '../../models/screener';
import { ScreenerResult } from '../../models/screener-result';
import { StateService } from '../../providers/state/state.service';
import { ScreenerService } from '../../providers/http/screener.service';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../../providers/http/trade.service';

@Component({
  selector: 'app-screener',
  templateUrl: './home.component.html'
})
export class ScreenerComponent implements OnInit 
{
  screeners: Screener[] = [];
  destory: Subject<boolean> = new Subject<boolean>();

  //
  // Constructor....
  //
  constructor(private stateService: StateService, private screenerService: ScreenerService, private tradeService: TradeService, private router: Router) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Default start timer
    let startTimer: number = (1000 * 10);

    // Get Data from cache
    this.screeners = this.stateService.GetScreens();  

    // Load page data.
    if(! this.screeners) 
    {
      startTimer = (1000 * 60);
      this.getScreeners();
    }

    // Reload the data every 1min after a 1 min delay to start
    Observable.timer(startTimer, (1000 * 60)).takeUntil(this.destory).subscribe(x => { this.getScreeners(); });    
  }

  //
  // OnDestroy
  //
  ngOnDestroy() 
  {
    this.destory.next();
    this.destory.complete();
  }  

  //
  // Get screeners.
  //
  getScreeners()
  {
    // Make api call to get screeners
    this.screenerService.get().subscribe(

      (res) => {
        this.screeners = res;

        // Make sure we have at least one screener
        if(this.screeners.length <= 0)
        {
          this.router.navigate(['/screener/add']);
          return;
        }

        // Load results.
        for (let i = 0; i < this.screeners.length; i++) 
        {
          this.screenerService.getResults(this.screeners[i].Id).subscribe((res) => {
            this.screeners[i].Results = res;
          });
        }

        // Store in site cache.
        this.stateService.SetScreens(this.screeners);
      }, 

      // Error
      (err: HttpErrorResponse) => {

        if (err.error instanceof Error) 
        {
          // A client-side or network error occurred. Handle it accordingly.
          console.log('An error occurred:', err.error);
        } else 
        {
          if (err.error.error == 'No Record Found.')
          {
            this.router.navigate(['/screener/add'], { queryParams: { action: 'first-run' } });
          }
        }
      }
    );
  }

  //
  // Get Spread string of a found result.
  //
  getSpread(result: ScreenerResult) : String
  {
    let exp = [];

    if (result.Legs.length <= 0) 
    {
      return null;
    }

    for (let i = 0; i < result.Legs.length; i++)
    {
      exp.push(result.Legs[i].OptionStrike);
    }

    return exp.join("/");
  }

  //
  // Get the expire string for a found result.
  //
  getExpire(result: ScreenerResult) : Date 
  {
    if (result.Legs.length <= 0)
    {
      return null;
    }

    return result.Legs[0].OptionExpire;
  }

  //
  // Place trade from result
  //
  trade(screen: Screener, result: ScreenerResult) {

    // Just a double check
    if(result.Legs.length <= 0)
    {
      return;
    }

    // Set values
    let tradeDetails = new TradeDetails();
    tradeDetails.Symbol = result.Legs[0].OptionUnderlying;
    tradeDetails.Class = "multileg";
    
    if(result.Credit > 0) 
    {
      tradeDetails.OrderType = "credit";
    } else 
    {
      tradeDetails.OrderType = "debit";
    }

    // TODO: Add configuration around this.
    tradeDetails.Duration = "gtc";
    
    // TODO: Add configuration around this. (or just selectors midpoint vs ask).
    tradeDetails.Price = result.MidPoint;

    // Build legs
    tradeDetails.Legs = [];

    for (let i = 0; i < result.Legs.length; i++) 
    {
      let side = "sell_to_close";
      
      // TODO: Get this from settings.
      let qty = 11;

      // TODO this will need work based on the type of screener.
      if (i == 1) 
      {
        side = "sell_to_open";
      } else 
      {
        side = "buy_to_open"
      }

      tradeDetails.Legs.push(new TradeOptionLegs().createNew(result.Legs[i], result.Legs[i].OptionExpire, result.Legs[i].OptionType, result.Legs[i].OptionStrike, side, qty));
    }

    // Open builder to place trade.
    this.tradeService.tradeEvent.emit(new TradeEvent().createNew("toggle-trade-builder", tradeDetails));


  }

}

/* End File */