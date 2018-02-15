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
  pageSize: number = environment.default_page_size;
  @Input() count: number;
  @Input() page: number;  
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
    if(this.count < environment.default_page_size)
    {
      this.page--;
    }

    this.pageClick.emit(this.page);
  }  
}

/* End File */