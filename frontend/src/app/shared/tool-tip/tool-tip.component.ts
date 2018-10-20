//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ElementRef, Input, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-shared-tool-tip',
  templateUrl: './tool-tip.component.html',
  styleUrls: [],
  host: { '(document:click)': 'onDocClick($event)' }  
})

export class ToolTipComponent implements OnInit 
{
  // Passed vars
  @Input() show: boolean;
  @Input() title: string;
  @Input() message: string;

  //
  // Construct.
  //
  constructor(private _eref: ElementRef) { }

  //
  // Ng Init
  //
  ngOnInit() {}

  //
  // Click anywhere on the screen.
  //
  onDocClick(event) 
  {
    if(! this._eref.nativeElement.contains(event.target)) 
    {
      this.closeToolTip();
    }
  }

  //
  // Open tool tip.
  //
  openToolTip()
  {
    if(this.show) 
    {
      this.show = false;
    } else 
    {
      this.show = true;
    }
  }

  //
  // Tool tip close.
  //
  closeToolTip()
  {
    this.show = false;
  }  

}

/* End File */
