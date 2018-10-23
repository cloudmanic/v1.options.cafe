//
// Date: 9/7/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//


import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { StateService } from '../../../providers/state/state.service';
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
  editId: number;
  runText: string = "Run";
  searching: boolean = false;
  runFirst: boolean = true;
  nameError: boolean = false;
  symbolError: boolean = false;
  newLoad: boolean = true;
  showDeleteScreener: boolean = false;
  results: ScreenerResult[] = [];
  screen: Screener = new Screener();
  itemSetttings: ScreenerItemSettings[] = [];


  //
  // Construct.
  //
  constructor(private stateService: StateService, private router: Router, private route: ActivatedRoute, private screenerService: ScreenerService) { }

  //
  // Ng Init
  //
  ngOnInit() 
  {
    // Is this an edit action?
    this.editId = this.route.snapshot.params['id'];

    // Setup widths
    let widths: number[] = [];

    for (let i = 0.5; i <= 500; i = i + .5)
    {
      widths.push(i);
    }

    // Setup days
    let days: number[] = [];

    for (let i = 0; i <= 500; i++) {
      days.push(i);
    }

    // Set item Settings
    this.itemSetttings.push(new ScreenerItemSettings('Spread Width', 'spread-width', 'select-number', ['=', '>'], widths, [], 0));
    this.itemSetttings.push(new ScreenerItemSettings('Open Credit', 'open-credit', 'input-number', ['>', '<', '='], [], [], 0.1));
    this.itemSetttings.push(new ScreenerItemSettings('Days To Expire', 'days-to-expire', 'select-number', ['>', '<', '='], days, [], 1.0));
    this.itemSetttings.push(new ScreenerItemSettings('Short Strike % Away', 'short-strike-percent-away', 'input-number', ['=', '>'], [], [], 0.1));

    // Default Values
    if(! this.editId) 
    {
      this.screen.Name = '';
      this.screen.Symbol = 'spy';
      this.screen.Strategy = 'put-credit-spread';

      // Screen
      this.screen.Items = [];
      this.screen.Items.push(new ScreenerItem(0, 0, 'spread-width', '=', '', 2.0, this.itemSetttings[0]));
      // this.screen.Items.push(new ScreenerItem(0, 0, 'min-credit', '=', '', 0.18, this.itemSetttings[1]));
      // this.screen.Items.push(new ScreenerItem(0, 0, 'days-to-expire', '<', '', 46, this.itemSetttings[2]));
      // this.screen.Items.push(new ScreenerItem(0, 0, 'short-strike-percent-away', '>', '', 4.0, this.itemSetttings[3]));

      // Run screen on load
      //this.runScreen();      
    } else 
    {
      // Make AJAX call to get this screener by ID.
      this.screenerService.getById(this.editId).subscribe((res) => {
        this.screen = res;

        for(let i = 0; i < this.screen.Items.length; i++)
        {
          for (let r = 0; r < this.itemSetttings.length; r++)
          {
            if (this.screen.Items[i].Key == this.itemSetttings[r].Key)
            {
              this.screen.Items[i].Settings = this.itemSetttings[r];
            }
          } 
        }

        this.runScreen();
      });
    }

  }

  //
  // Add 
  // 
  addFilter()
  {
    this.screen.Items.push(new ScreenerItem(0, 0, '', '=', '', 2.0, this.itemSetttings[0]));    
  }

  //
  // Filter change.
  //
  filterChange()
  {
    this.runFirst = true;
    this.nameError = false;
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
    this.newLoad = false;
    this.nameError = false;
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

    // Must have a name for the screener
    if(this.screen.Name.length <= 0)
    {
      this.nameError = true;    
      return false;
    } else
    {
      this.nameError = false;
    }

    this.runFirst = true;

    // Send API call to server to get the results for this screen.
    if(this.editId) 
    {
      // Update screen.
      this.screenerService.submitUpdate(this.screen).subscribe((res) => {
        this.stateService.SetScreens(null);
        this.router.navigate(['/screener'], { queryParams: { update: this.screen.Id, title: this.screen.Name } });
      });
    } else 
    {
      // Create new screen.
      this.screenerService.submitScreen(this.screen).subscribe((res) => {
        this.stateService.SetScreens(null);
        this.router.navigate(['/screener'], { queryParams: { new: res.Id } });
      });
    }

    return true;
  }

  //
  // Delete Screen
  //
  deleteScreen()
  {
    this.showDeleteScreener = true;
  }

  //
  // On Delete Screener
  //
  onDeleteScreen()
  {
    // Send API call to server to delete this screen.
    this.screenerService.deleteById(this.screen.Id).subscribe((res) => {
      this.showDeleteScreener = false;
      this.stateService.SetScreens(null);
      this.router.navigate(['/screener'], { queryParams: { deleted: this.screen.Name } });
    });
  }
}

/* End File */