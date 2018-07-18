//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Rx';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit } from '@angular/core';
import { Screener } from '../../models/screener';
import { StateService } from '../../providers/state/state.service';
import { ScreenerService } from '../../providers/http/screener.service';

@Component({
  selector: 'app-screener',
  templateUrl: './home.component.html'
})
export class ScreenerComponent implements OnInit 
{
  screeners: Screener[]
  destory: Subject<boolean> = new Subject<boolean>(); 

  //
  // Constructor....
  //
  constructor(private stateService: StateService, private screenerService: ScreenerService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Load page data.
    this.getData();

    // Reload the data every 1min after a 1 min delay to start
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
    this.getScreeners();
  }

  //
  // Get screeners.
  //
  getScreeners()
  {
    // Make api call to get screeners
    this.screenerService.get().subscribe((res) => {
      this.screeners = res;

      // Load results.
      for (let i = 0; i < this.screeners.length; i++) 
      {
        this.screenerService.getResults(this.screeners[i].Id).subscribe((res) => {
          this.screeners[i].Results = res;

          console.log(this.screeners[i]);
        });
      }

      //this.stateService.SetDashboardRank(this.rank);
    });
  }

}

/* End File */