//
// Date: 9/7/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//


import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Settings } from '../../../models/settings';
import { StateService } from '../../../providers/state/state.service';
import { ScreenerService } from '../../../providers/http/screener.service';
import { Screener, ScreenerItem, ScreenerItemSettings } from '../../../models/screener';
import { ScreenerResult } from '../../../models/screener-result';
import { SettingsService } from '../../../providers/http/settings.service';
import { DropdownAction } from '../../../shared/dropdown-select/dropdown-select.component';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../../../providers/http/trade.service';

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
  settings: Settings = new Settings();

  actions: DropdownAction[] = null;     

  itemSetttings: Map<string, ScreenerItemSettings[]> = new Map<string, ScreenerItemSettings[]>();

  helperPost: string = "For more information checkout <a target=\"_blank\" href=\"https://cloudmanic.groovehq.com/knowledge_base/topics/using-the-options-screener\">The Options Screener</a>.";
  strategyHelpTitle: string = "Options Strategy";
  strategyHelpBody: string = "Here you can select what type of options strategy you want to screen for. An options strategy is a trade that is betting the underlying equity moves in a certain way over the course of a define period of time. For more information on possible strategies checkout this <a href=\"https://www.investopedia.com/articles/active-trading/040915/guide-option-trading-strategies-beginners.asp\" target=\"_blank\">post</a>. " + this.helperPost;

  underlyingSymbolHelpTitle: string = "Underlying Symbol";
  underlyingSymbolHelpBody: string = "Options Cafe screeners are limited to screening for only one underlying symbol (or equity) at a time. If you would like to keep an eye on more than one underlying symbols simply create and save more than one screener. " + this.helperPost;

  toolTips: Map<string, ToolTipItem> = new Map<string, ToolTipItem>();

  //
  // Construct.
  //
  constructor(private tradeService: TradeService, private stateService: StateService, private router: Router, private route: ActivatedRoute, private screenerService: ScreenerService, private settingsService: SettingsService) 
  { 
    this.setupToolTips();
    this.setupItemSettings();
    this.settings = this.stateService.GetSettings();
  }

  //
  // Ng Init
  //
  ngOnInit() 
  {
    // Setup Dropdown actions
    this.setupDropdownActions();

    // Load settings data
    this.loadSettingsData();    

    // Is this an edit action?
    this.editId = this.route.snapshot.params['id'];

    // Default Values
    if(! this.editId) 
    {
      this.screen.Name = '';
      this.screen.Symbol = 'spy';
      this.screen.Strategy = 'put-credit-spread';

      // Screen
      this.screen.Items = [];
      //this.screen.Items.push(new ScreenerItem(0, 0, 'spread-width', '=', '', 2.0));
      // this.screen.Items.push(new ScreenerItem(0, 0, 'min-credit', '=', '', 0.18));
      // this.screen.Items.push(new ScreenerItem(0, 0, 'days-to-expire', '<', '', 46));
      // this.screen.Items.push(new ScreenerItem(0, 0, 'short-strike-percent-away', '>', '', 4.0);

      // Run screen on load
      //this.runScreen();      
    } else 
    {
      // Make AJAX call to get this screener by ID.
      this.screenerService.getById(this.editId).subscribe((res) => {
        this.screen = res;

        for(let i = 0; i < this.screen.Items.length; i++)
        {
          for(let r = 0; r < this.itemSetttings[this.screen.Strategy].length; r++)
          {
            if(this.screen.Items[i].Key == this.itemSetttings[this.screen.Strategy][r].Key)
            {
              this.screen.Items[i].Settings = this.itemSetttings[this.screen.Strategy][r];
            }
          } 
        }

        this.runScreen();
      });
    }

  }

  //
  // Load settings data.
  //
  loadSettingsData() 
  {
    this.settingsService.get().subscribe((res) => {
      this.settings = res;

      console.log(this.settings);
      this.stateService.SetSettings(res);
    });
  } 

  //
  // Setup Drop down actions.
  //
  setupDropdownActions() {
    let das = []

    // First action
    let da1 = new DropdownAction();
    da1.title = "Place Trade";

    // Place order action
    da1.click = (result: ScreenerResult) => {
      this.trade(this.screen, result);

    };

    das.push(da1);

    this.actions = das;
  }   

  //
  // Place trade from result
  //
  trade(screen: Screener, result: ScreenerResult) {

    // Just a double check
    if (result.Legs.length <= 0) 
    {
      return;
    }

    // Set values
    let tradeDetails = new TradeDetails();
    tradeDetails.Symbol = result.Legs[0].OptionUnderlying;
    tradeDetails.Class = "multileg";

    if (result.Credit > 0) 
    {
      tradeDetails.OrderType = "credit";
    } else 
    {
      tradeDetails.OrderType = "debit";
    }

    // TODO: Add configuration around this.
    tradeDetails.Duration = "gtc";

    // Default Price
    tradeDetails.Price = result.MidPoint;

    // Build legs
    tradeDetails.Legs = [];

    // Figure out side based on strategy
    let qty = 1;
    let sides: string[] = ["buy_to_open", "sell_to_open", "sell_to_open", "buy_to_open"];

    switch (screen.Strategy) 
    {
      case "reverse-iron-condor":
        sides = ["sell_to_open", "buy_to_open", "buy_to_open", "sell_to_open"];
      break;

      case "put-credit-spread":
        qty = this.settings.StrategyPcsLots;
        sides = ["buy_to_open", "sell_to_open", "sell_to_open", "buy_to_open"];

        // TODO: Make this work. 

        // if(this.settings.StrategyPcsOpenPrice == "mid-point") 
        // {
        //   tradeDetails.Price = result.MidPoint;
        // }

        // if(this.settings.StrategyPcsOpenPrice == "ask") 
        // {
        //   tradeDetails.Price = result.MidPoint;
        // }        
      break;     
    }

    for (let i = 0; i < result.Legs.length; i++) 
    {
      let side: string = sides[i];
      tradeDetails.Legs.push(new TradeOptionLegs().createNew(result.Legs[i], result.Legs[i].OptionExpire, result.Legs[i].OptionType, result.Legs[i].OptionStrike, side, qty));
    }

    // Open builder to place trade.
    this.tradeService.tradeEvent.emit(new TradeEvent().createNew("toggle-trade-builder", tradeDetails));


  }  

  //
  // Add 
  // 
  addFilter()
  {
    let newRow: ScreenerItem = new ScreenerItem(0, 0, 'spread-width', '>', '', 2.0);
    newRow.Settings = this.itemSetttings[this.screen.Strategy][0];
    newRow.Operator = newRow.Settings.Operators[0];
    this.screen.Items.push(newRow);    
  }

  //
  // Strategy Change
  //
  strategyChange()
  {
    this.screen.Items = [];
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

  //
  // Setup tool tips.
  //
  setupToolTips()
  {
    // Add the helper text to the tool tip map
    this.toolTips.set("open-debit", new ToolTipItem("Open Debit", "This is how much you are willing to pay to open this trade. " + this.helperPost));
    this.toolTips.set("put-leg-width", new ToolTipItem("Put Leg Width", "Set how wide you want your put leg to spread. " + this.helperPost));
    this.toolTips.set("call-leg-width", new ToolTipItem("Call Leg Width", "Set how wide you want your call leg to spread. " + this.helperPost));   
    this.toolTips.set("put-leg-percent-away", new ToolTipItem("Put Leg Percent Away", "How far away from the current stock price should your first put leg should be. This is a percentage. " + this.helperPost));
    this.toolTips.set("call-leg-percent-away", new ToolTipItem("Call Leg Percent Away", "How far away from the current stock price should your first call leg should be. This is a percentage. " + this.helperPost));

    this.toolTips.set("spread-width", new ToolTipItem("Spread Width", "Set how wide you want your spreads to be. For example with a put credit spread where the short strike is 200 and the long strike is 195 this spread width would be 5. " + this.helperPost));
    this.toolTips.set("open-credit", new ToolTipItem("Open Credit", "In most cases we are only interested in spreads that give at least a set minimum for the opening credit. For example if you wanted to open a put credit spread with an opening credit of $0.15 you would receive $15 for each lot. " + this.helperPost));
    this.toolTips.set("days-to-expire", new ToolTipItem("Days To Expire", "Days to expire is the number of days until the trade would expire. Often times we are not interested in trades that expire too far into the future or trades that expire too soon. " + this.helperPost));
    this.toolTips.set("short-strike-percent-away", new ToolTipItem("Short Strike % Away", "To help define our risk we might only be interested in trades where the current underlying stock is some percentage away from your short options leg. For example if we are trading put credit spreads on the SPY and the SPY is currently trading at $100. We could set this value to 5.0 and we would only be presented with trades where the short leg strike price is $95 or less. " + this.helperPost));  
  }

  //
  // Setup item settings
  //
  setupItemSettings()
  {
    // Setup widths
    let widths: number[] = [];

    for (let i = 0.5; i <= 500; i = i + .5) 
    {
      widths.push(i);
    }

    // Setup days
    let days: number[] = [];

    for (let i = 0; i <= 500; i++) 
    {
      days.push(i);
    }    

    // Set item Settings: put-credit-spread
    let pcs: ScreenerItemSettings[] = [];
    pcs.push(new ScreenerItemSettings('Spread Width', 'spread-width', 'select-number', ['='], widths, [], 0));
    pcs.push(new ScreenerItemSettings('Open Credit', 'open-credit', 'input-number', ['>', '<', '='], [], [], 0.1));
    pcs.push(new ScreenerItemSettings('Days To Expire', 'days-to-expire', 'select-number', ['<', '>', '='], days, [], 30));
    pcs.push(new ScreenerItemSettings('Short Strike % Away', 'short-strike-percent-away', 'input-number', ['>', '<'], [], [], 4.0));
    this.itemSetttings["put-credit-spread"] = pcs;

    // Set item Settings: reverse-iron-condor
    let ric: ScreenerItemSettings[] = [];
    ric.push(new ScreenerItemSettings('Days To Expire', 'days-to-expire', 'select-number', ['<', '>', '='], days, [], 30));
    ric.push(new ScreenerItemSettings('Open Debit', 'open-debit', 'input-number', ['<', '>', '='], [], [], 1.00));
    ric.push(new ScreenerItemSettings('Put Leg Width', 'put-leg-width', 'select-number', ['='], widths, [], 2.00));
    ric.push(new ScreenerItemSettings('Call Leg Width', 'call-leg-width', 'select-number', ['='], widths, [], 2.00));
    ric.push(new ScreenerItemSettings('Put Leg % Away', 'put-leg-percent-away', 'input-number', ['>', '<'], [], [], 4.0));
    ric.push(new ScreenerItemSettings('Call Leg % Away', 'call-leg-percent-away', 'input-number', ['>', '<'], [], [], 4.0));
    this.itemSetttings["reverse-iron-condor"] = ric;
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