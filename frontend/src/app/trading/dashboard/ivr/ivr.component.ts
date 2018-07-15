//
// Date: 2/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Rx';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit } from '@angular/core';
import { Rank } from '../../../models/rank';
import { StateService } from '../../../providers/state/state.service';
import { QuotesService } from '../../../providers/http/quotes.service';

@Component({
  selector: '[app-trading-ivr]',
  templateUrl: './ivr.component.html'
})

export class IvrComponent implements OnInit 
{
  rank: Rank = null;
  destory: Subject<boolean> = new Subject<boolean>();

  //
  // Constructor....
  //
  constructor(private stateService: StateService, private quotesService: QuotesService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Set rank from cache
    this.rank = this.stateService.GetDashboardRank();

    // Get data
    this.getData();

    // Reload the ivr data every 1min after a 1 min delay to start
    Observable.timer((1000 * 60), (1000 * 60)).takeUntil(this.destory).subscribe(x => { this.getData(); });
  } 

  //
  // OnDestroy
  //
  ngOnDestroy() {
    this.destory.next();
    this.destory.complete();
  }     

  //
  // Get data
  //
  getData() 
  {
    // Make api call to get rank data.
    this.quotesService.getSymbolRank("vix").subscribe((res) => {
      this.rank = res;
      this.stateService.SetDashboardRank(this.rank);
    });
  }

}

/* End File */