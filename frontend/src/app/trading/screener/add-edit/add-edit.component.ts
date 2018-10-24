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

  helperPost: string = "For more information checkout <a target=\"_blank\" href=\"https://cloudmanic.groovehq.com/knowledge_base/topics/using-the-options-screener\">The Options Screener</a>.";
  strategyHelpTitle: string = "Options Strategy";
  strategyHelpBody: string = "Here you can select what type of options strategy you want to screen for. An options strategy is a trade that is betting the underlying equity moves in a certain way over the course of a define period of time. For more information on possible strategies checkout this <a href=\"https://www.investopedia.com/articles/active-trading/040915/guide-option-trading-strategies-beginners.asp\" target=\"_blank\">post</a>. " + this.helperPost;

  underlyingSymbolHelpTitle: string = "Underlying Symbol";
  underlyingSymbolHelpBody: string = "Options Cafe screeners are limited to screening for only one underlying symbol (or equity) at a time. If you would like to keep an eye on more than one underlying symbols simply create and save more than one screener. " + this.helperPost;

  toolTips: Map<string, ToolTipItem> = new Map<string, ToolTipItem>();

  //
  // Construct.
  //
  constructor(private stateService: StateService, private router: Router, private route: ActivatedRoute, private screenerService: ScreenerService) 
  { 
    // Add the helper text to the tool tip map
    this.toolTips.set("spread-width", new ToolTipItem("Spread Width", "Set how wide you want your spreads to be. For example with a put credit spread where the short strike is 200 and the long strike is 195 this spread width would be 5. " + this.helperPost));
    this.toolTips.set("open-credit", new ToolTipItem("Open Credit", "In most cases we are only interested in spreads that give at least a set minimum for the opening credit. For example if you wanted to open a put credit spread with an opening credit of $0.15 you would receive $15 for each lot. " + this.helperPost));
    this.toolTips.set("days-to-expire", new ToolTipItem("Days To Expire", "Days to expire is the number of days until the trade would expire. Often times we are not interested in trades that expire too far into the future or trades that expire too soon. " + this.helperPost));
    this.toolTips.set("short-strike-percent-away", new ToolTipItem("Short Strike % Away", "To help define our risk we might only be interested in trades where the current underlying stock is some percentage away from your short options leg. For example if we are trading put credit spreads on the SPY and the SPY is currently trading at $100. We could set this value to 5.0 and we would only be presented with trades where the short leg strike price is $95 or less. " + this.helperPost));            
  }

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
    this.itemSetttings.push(new ScreenerItemSettings('Short Strike % Away', 'short-strike-percent-away', 'input-number', ['>', '<'], [], [], 0.1));

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
    this.screen.Items.push(new ScreenerItem(0, 0, 'open-credit', '>', '', 2.0, this.itemSetttings[0]));    
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

  //
  // Show helper tool tips
  //
  openToolTip(key: string)
  {
    if(this.toolTips.get(key).Show) 
    {
      this.toolTips.get(key).Show = false;
    } else
    {
      this.toolTips.get(key).Show = true;
    }
  }
}

//
// Tool tip Item
//
export class ToolTipItem 
{
  Title: string = "";
  Body: string = "";
  Show: boolean = false;

  //
  // Construct.
  //
  constructor(title: string, body: string)
  {
    this.Title = title;
    this.Body = body;
  }
}

/* End File */