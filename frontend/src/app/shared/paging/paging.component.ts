//
// Date: 2/14/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//
//

import { environment } from '../../../environments/environment';
import { Component, OnInit, Input, Output, OnChanges, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-paging',
  templateUrl: './paging.component.html'
})

export class PagingComponent implements OnInit {
  @Input() count: number;
  @Input() page: number;
  @Input() limit: number;
  @Input() noLimitCount: number;  
  @Output('pageClick') pageClick: EventEmitter<number> = new EventEmitter<number>();

  //
  // Construct
  //
  constructor() { }

  //
  // OnInit
  //
  ngOnInit() { }

  //
  // On paging click - Prev
  //
  onPrev()
  {
    this.page--;

    // Are we at the start?
    if(this.page < 1)
    {
      this.page = 1;
    }

    this.pageClick.emit(this.page);
  }

  //
  // On paging click - Next
  //
  onNext()
  {
    this.page++;    

    // Did we go too far with Next?
    if(this.count < this.limit)
    {
      this.page--;
    }

    this.pageClick.emit(this.page);
  }  

  //
  // Get range start
  //
  getRangeStart() : number {
    if(this.count == 0)
    {
      return Number(this.count);
    }

    return (this.page * this.limit) - this.limit + 1;
  }

  //
  // Get range end
  //
  getRangeEnd() : number {
    if(this.count == 0)
    {
      return Number(this.count);
    }
    
    return this.getRangeStart() + this.count - 1;
  }  
}

/* End File */