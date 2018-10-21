//
// Date: 6/28/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input, ElementRef } from '@angular/core';

@Component({
  selector: 'app-dropdown-select',
  templateUrl: './dropdown-select.component.html',
  host: { '(document:click)': 'onDocClick($event)' }
})

export class DropdownSelectComponent implements OnInit 
{
  @Input() data: any;
  @Input() actions: DropdownAction[];
  active: boolean = false;

  //
  // Constructor
  //
  constructor(private _eref: ElementRef) { }

  //
  // NgInit
  //
  ngOnInit() {}

  //
  // Click anywhere on the screen.
  //
  onDocClick(event) {

    // Remove active buttons
    if (!this._eref.nativeElement.contains(event.target)) 
    {
      this.active = false;
    }
  }

  //
  // Toggle drop down
  //
  onToggleDropDown()
  {
    if(this.active)
    {
      this.active = false;
    } else
    {
      this.active = true;
    } 
  }
}

export class DropdownAction 
{
  title: string = "";
  click: Function = null;
  section: boolean = false;
}

/* End File */
