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
    // Set item Settings
    this.itemSetttings.push(new ScreenerItemSettings('Spread Width', 'select', ['=', '>', '>='], ['0.5', '1.0', '1.5', '2.0', '2.5', '3.0']));
    this.itemSetttings.push(new ScreenerItemSettings('Open Credit', 'input', ['=', '>', '>='], []));
    this.itemSetttings.push(new ScreenerItemSettings('Max Days To Expire', 'input', ['='], []));
    this.itemSetttings.push(new ScreenerItemSettings('Short Strike Percent Away', 'input', ['>', '>='], []));


    // Screen
    this.screen.Items = [];
    this.screen.Items.push(new ScreenerItem(0, 0, '', '=', '', 0, this.itemSetttings[0]));


    //console.log(this.itemSetttings);
  }

}

/* End File */