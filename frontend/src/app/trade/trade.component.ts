//
// Date: 4/2/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Symbol } from '../models/symbol';
import { Component, OnInit } from '@angular/core';
import { OptionsChainService } from '../providers/http/options-chain.service';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../providers/http/trade.service';
import { WebsocketService } from '../providers/http/websocket.service';


@Component({
  selector: 'app-trade',
  templateUrl: './trade.component.html'
})

export class TradeComponent implements OnInit 
{
  quotes = {}
  symbol: Symbol;
  typeAheadSymbol: Symbol;
  symbolExpirations: Date[] = [];
  tradeDetails: TradeDetails = new TradeDetails();
  showTradeBuilder: boolean = true;

  //
  // Construct.
  //
  constructor(private websocketService: WebsocketService, private tradeService: TradeService, private optionsChainService: OptionsChainService) 
  { 
    // // Set Defaults (also used for development)
    // this.tradeDetails.Symbol = "SPY";
    // this.tradeDetails.Class = "multileg";
    // this.tradeDetails.OrderType = "credit";
    // this.tradeDetails.Duration = "gtc";
    // this.tradeDetails.Price = 0.21;
    
    // // Build legs
    // this.tradeDetails.Legs = [];
    // this.tradeDetails.Legs.push(new TradeOptionLegs().createNew("SPY180402P00262000", moment("2018-05-04").toDate(), "Puts", 190, "sell_to_open", 10));
    // this.tradeDetails.Legs.push(new TradeOptionLegs().createNew("SPY180402P00264000", moment("2018-05-04").toDate(), "Calls", 190, "buy_to_open", 10));   
    

    // // Load data (we call this only when we pass in a full trade)
    // this.loadChainForAllLegs();

    //console.log(this.tradeDetails);


    // Subscribe to data updates from the quotes - Market Quotes
    this.websocketService.quotePushData.subscribe(data => {
      this.quotes[data.symbol] = data;
    }); 
  }

  //
  // OnInit.
  //
  ngOnInit() 
  {
    // Subscribe trade events
    this.tradeService.tradeEvent.subscribe(data => {
      this.manageTradeEvent(data);
    });
  }

  //
  // Manage new trade events
  //
  manageTradeEvent(data: TradeEvent)
  {
    // Populate the form.
    this.tradeDetails = data.TradeDetails;

    // Set default if need be.
    if(this.tradeDetails.Class == undefined) 
    {
      this.tradeDetails.Class = "equity";
    }

    // Toggle form.
    this.toggleTradeBuilder();    
  }

  //
  // Toggle the trade builder
  //
  toggleTradeBuilder()
  {
    if(this.showTradeBuilder)
    {
      this.showTradeBuilder = false;
    } else
    {
      this.showTradeBuilder = true;      
    }
  }

  //
  // Set trade class.
  //
  setTradeClass(type: string)
  {
    this.tradeDetails.Class = type;
  }

  //
  // onSearchTypeAheadClick() 
  //
  onSearchTypeAheadClick(symbol: Symbol) 
  {
    if (typeof symbol == "undefined") {
      return;
    }
    
    // Set the symbol
    this.typeAheadSymbol = symbol;
  } 

  //
  // onSearchSubmit()
  //
  onSearchSubmit() 
  {
    // Set the symbol
    this.symbol = this.typeAheadSymbol;
    this.tradeDetails.Symbol = this.symbol.ShortName;

    // Load expire dates
    this.loadExpireDates(true);    
  }

  //
  // Load expire dates
  //
  loadExpireDates(addLeg: boolean)
  {
    this.symbolExpirations = [];
    
    // Make API call to get option expire dates.
    this.tradeService.getOptionExpirations(this.tradeDetails.Symbol).subscribe(data => {
      for (let i = 0; i < data.length; i++)
      {
        this.symbolExpirations.push(data[i]);
      }

      // Add first leg
      if(addLeg) 
      {
        this.tradeDetails.Legs = [];
        this.tradeDetails.Legs.push(new TradeOptionLegs());
        this.tradeDetails.Legs[0].Type = "Puts";
        this.tradeDetails.Legs[0].Side = "buy_to_open";
        this.tradeDetails.Legs[0].Expire = this.symbolExpirations[0];
        this.onExpireChange(this.tradeDetails.Legs[0], 0);
      }
    });
  }

  //
  // Load chain for all legs (normally called at the start when we pass in a full trade)
  //
  loadChainForAllLegs()
  {
    for (let i = 0; i < this.tradeDetails.Legs.length; i++)
    {
      this.onExpireChange(this.tradeDetails.Legs[i], i);
    }
  }

  //
  // On expire change.
  //
  onExpireChange(leg: TradeOptionLegs, index: number) 
  {
    // Api call to get the strikes for this chain.
    this.tradeService.getOptionStrikesBySymbolExpiration(this.tradeDetails.Symbol, leg.Expire).subscribe(data => {
      this.tradeDetails.Legs[index].Strikes = data;

      // Make sure our strike is still valid.
      let found = false;

      for(let i = 0; i < this.tradeDetails.Legs[index].Strikes.length; i++)
      {
        if(this.tradeDetails.Legs[index].Strikes[i] == this.tradeDetails.Legs[index].Strike)
        {
          found = true;
        }
      }

      if(! found)
      {
        this.tradeDetails.Legs[index].Strike = this.tradeDetails.Legs[index].Strikes[0];
      }
    });

    // // Api call to get the option chain.
    // this.optionsChainService.getChainBySymbolExpire(this.tradeDetails.Symbol, leg.Expire).subscribe(data => {
    //   this.tradeDetails.Legs[index].Chain = data;
    // });  
  }

  //
  // onTypeChange - When we pick a call or a put
  //
  onTypeChange(leg: TradeOptionLegs, index: number)
  {
    //console.log(this.tradeDetails);
  }

  //
  // Duplicate leg
  //
  duplicateLeg(leg: TradeOptionLegs)
  {
    let newLeg = new TradeOptionLegs();
    newLeg.Symbol = leg.Symbol;
    newLeg.Expire = leg.Expire;
    newLeg.Type = leg.Type;
    newLeg.Strike = leg.Strike;
    newLeg.Side = leg.Side;
    newLeg.Qty = leg.Qty;
    newLeg.Strikes = leg.Strikes;
    this.tradeDetails.Legs.push(newLeg); 
  }

  //
  // Remove leg
  //
  removeLeg(leg: TradeOptionLegs, index: number) {
    this.tradeDetails.Legs.splice(index, 1);
  }  
}

/* End File */