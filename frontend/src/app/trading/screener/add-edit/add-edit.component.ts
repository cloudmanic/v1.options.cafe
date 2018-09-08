//
// Date: 9/7/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//


import { Component, OnInit } from '@angular/core';
import { Screener, ScreenerItem, ScreenerItemSettings } from '../../../models/screener';
import { ScreenerResult } from '../../../models/screener-result';

@Component({
  selector: 'app-add-edit',
  templateUrl: './add-edit.component.html',
  styleUrls: []
})

export class AddEditComponent implements OnInit 
{
  runFirst: boolean = true;
  symbolError: boolean = false;
  screen: Screener = new Screener();
  itemSetttings: ScreenerItemSettings[] = [];


  //
  // Construct.
  //
  constructor() { }

  //
  // Ng Init
  //
  ngOnInit() 
  {
    // Default Values
    this.screen.Name = '';
    this.screen.Symbol = '';
    this.screen.Strategy = 'put-credit-spread';

    // Setup widths
    let widths: number[] = [];

    for (let i = 0.5; i <= 500; i = i + .5)
    {
      widths.push(i);
    }

    // Setup days
    let days: number[] = [];

    for (let i = 1; i <= 500; i++) {
      days.push(i);
    }

    // Set item Settings
    this.itemSetttings.push(new ScreenerItemSettings('Spread Width', 'select-number', ['=', '>'], widths, [], 0));
    this.itemSetttings.push(new ScreenerItemSettings('Open Credit', 'input-number', ['=', '>'], [], [], 0.1));
    this.itemSetttings.push(new ScreenerItemSettings('Max Days To Expire', 'select-number', ['='], days, [], 1.0));
    this.itemSetttings.push(new ScreenerItemSettings('Short Strike % Away', 'input-number', ['>'], [], [], 0.1));


    // Screen
    this.screen.Items = [];
    this.screen.Items.push(new ScreenerItem(0, 0, '', '=', '', 2.0, this.itemSetttings[0]));
  }

  //
  // Add 
  // 
  addFilter()
  {
    this.screen.Items.push(new ScreenerItem(0, 0, '', '=', '', 2.0, this.itemSetttings[0]));    
  }

  //
  // Validate screen
  //
  validateScreen() : boolean
  {
    console.log(this.screen);

    // Validate symbol
    if(this.screen.Symbol.length <= 0)
    {
      this.symbolError = true;
      return false;
    } else
    {
      this.symbolError = false;
    }

    return true;
  }

  //
  // Run Screen
  //
  runScreen() : boolean
  {
    // First we validate
    if(! this.validateScreen())
    {
      return false;
    }

    this.runFirst = false;

    return true;
  }

  //
  // Save Screen
  //
  saveScreen(): boolean {

    // First we validate
    if(! this.validateScreen()) 
    {
      return false;
    }

    console.log(this.screen);

    this.runFirst = true;

    return true;
  }
}

/* End File */