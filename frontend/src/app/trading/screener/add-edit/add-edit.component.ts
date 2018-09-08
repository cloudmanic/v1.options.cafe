//
// Date: 9/7/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//


import { Component, OnInit } from '@angular/core';
import { ScreenerService } from '../../../providers/http/screener.service';
import { Screener, ScreenerItem, ScreenerItemSettings } from '../../../models/screener';
import { ScreenerResult } from '../../../models/screener-result';

@Component({
  selector: 'app-add-edit',
  templateUrl: './add-edit.component.html',
  styleUrls: []
})

export class AddEditComponent implements OnInit 
{
  runText: string = "Run";
  searching: boolean = false;
  runFirst: boolean = true;
  symbolError: boolean = false;
  results: ScreenerResult[] = [];
  screen: Screener = new Screener();
  itemSetttings: ScreenerItemSettings[] = [];


  //
  // Construct.
  //
  constructor(private screenerService: ScreenerService) { }

  //
  // Ng Init
  //
  ngOnInit() 
  {
    // Default Values
    this.screen.Name = '';
    this.screen.Symbol = 'spy';
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
    this.itemSetttings.push(new ScreenerItemSettings('Spread Width', 'spread-width', 'select-number', ['=', '>'], widths, [], 0));
    this.itemSetttings.push(new ScreenerItemSettings('Open Credit', 'min-credit', 'input-number', ['=', '>'], [], [], 0.1));
    this.itemSetttings.push(new ScreenerItemSettings('Max Days To Expire', 'max-days-to-expire', 'select-number', ['='], days, [], 1.0));
    this.itemSetttings.push(new ScreenerItemSettings('Short Strike % Away', 'short-strike-percent-away', 'input-number', ['>'], [], [], 0.1));


    // Screen
    this.screen.Items = [];
    this.screen.Items.push(new ScreenerItem(0, 0, 'spread-width', '=', '', 2.0, this.itemSetttings[0]));
    this.screen.Items.push(new ScreenerItem(0, 0, 'min-credit', '=', '', 0.18, this.itemSetttings[1]));
    this.screen.Items.push(new ScreenerItem(0, 0, 'max-days-to-expire', '=', '', 45, this.itemSetttings[2]));
    this.screen.Items.push(new ScreenerItem(0, 0, 'short-strike-percent-away', '>', '', 4.0, this.itemSetttings[3]));

    // Run screen on load
    this.runScreen();
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
    // Validate symbol
    if(this.screen.Symbol.length <= 0)
    {
      this.symbolError = true;
      return false;
    } else
    {
      this.symbolError = false;
    }

    // Set the keys
    for(let i = 0; i < this.screen.Items.length; i++)
    {
      this.screen.Items[i].Key = this.screen.Items[i].Settings.Key;
    }

    return true;
  }

  //
  // Run Screen
  //
  runScreen() : boolean
  {
    // If we are searching do nothing
    if(this.searching)
    {
      return false;
    }

    // First we validate
    if(! this.validateScreen())
    {
      return false;
    }

    // Set the state before searching.
    this.results = [];
    this.searching = true;
    this.runText = "Searching...";

    // Send API call to server to get the results for this screen.
    this.screenerService.submitScreenForResults(this.screen).subscribe((res) => {
      this.results = res;
      this.runText = "Run";
      this.runFirst = false;
      this.searching = false;
    });

    return true;
  }

  //
  // Save Screen
  //
  saveScreen(): boolean {

    // If we are searching do nothing
    if (this.searching) {
      return false;
    }

    // First we validate
    if(! this.validateScreen()) 
    {
      return false;
    }

    this.runFirst = true;

    return true;
  }
}

/* End File */