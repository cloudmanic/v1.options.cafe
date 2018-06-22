//
// Date: 6/21/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Symbol } from '../../models/symbol';
import { Component, OnInit, Input } from '@angular/core';
import { SymbolService } from '../../providers/http/symbol.service';
import { OptionsChainService } from '../../providers/http/options-chain.service';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../../providers/http/trade.service';

@Component({
  selector: 'app-trade-multi-leg',
  templateUrl: './trade-multi-leg.component.html'
})

export class TradeMultiLegComponent implements OnInit 
{
  @Input() quotes = {};
  @Input() tradeDetails: TradeDetails;

  symbol: Symbol;
  typeAheadSymbol: Symbol;
  symbolExpirations: Date[] = [];

  //
  // Construct.
  //
  constructor(private tradeService: TradeService, private optionsChainService: OptionsChainService, private symbolService: SymbolService) {}

  //
  // OnInit.
  //
  ngOnInit() {}

  //
  // Preview Trade
  //
  previewTrade() {
    alert('Preview');

    console.log(this.tradeDetails);
  }

  //
  // onSearchTypeAheadClick() 
  //
  onSearchTypeAheadClick(symbol: Symbol) {
    if (typeof symbol == "undefined") {
      return;
    }

    // Set the symbol
    this.typeAheadSymbol = symbol;
  }

  //
  // onSearchSubmit()
  //
  onSearchSubmit() {
    // Set the symbol
    this.symbol = this.typeAheadSymbol;
    this.tradeDetails.Symbol = this.symbol.ShortName;

    // Make sure we are getting quotes from the websocket for this symbol.
    this.symbolService.addActiveSymbol(this.symbol.ShortName).subscribe();

    // Load expire dates
    this.loadExpireDates(true);
  }

  //
  // Load expire dates
  //
  loadExpireDates(addLeg: boolean) {
    this.symbolExpirations = [];

    // Make API call to get option expire dates.
    this.tradeService.getOptionExpirations(this.tradeDetails.Symbol).subscribe(data => {
      for (let i = 0; i < data.length; i++) {
        this.symbolExpirations.push(data[i]);
      }

      // Add first leg
      if (addLeg) {
        this.tradeDetails.Legs = [];
        this.tradeDetails.Legs.push(new TradeOptionLegs());
        this.tradeDetails.Legs[0].Type = "Put";
        this.tradeDetails.Legs[0].Side = "buy_to_open";
        this.tradeDetails.Legs[0].Expire = this.symbolExpirations[0];
        this.onExpireChange(this.tradeDetails.Legs[0], 0);
      }
    });
  }

  //
  // Load chain for all legs (normally called at the start when we pass in a full trade)
  //
  loadChainForAllLegs() {
    for (let i = 0; i < this.tradeDetails.Legs.length; i++) {
      this.onExpireChange(this.tradeDetails.Legs[i], i);
      this.loadLegSymbol(this.tradeDetails.Legs[i]);
    }
  }

  //
  // On expire change.
  //
  onExpireChange(leg: TradeOptionLegs, index: number) {
    // Api call to get the strikes for this chain.
    this.tradeService.getOptionStrikesBySymbolExpiration(this.tradeDetails.Symbol, leg.Expire).subscribe(data => {
      this.tradeDetails.Legs[index].Strikes = data;

      // Make sure our strike is still valid.
      let found = false;

      for (let i = 0; i < this.tradeDetails.Legs[index].Strikes.length; i++) {
        if (this.tradeDetails.Legs[index].Strikes[i] == this.tradeDetails.Legs[index].Strike) {
          found = true;
        }
      }

      if (!found) {
        this.tradeDetails.Legs[index].Strike = this.tradeDetails.Legs[index].Strikes[0];
      }

      // Load leg quotes
      this.loadLegSymbol(leg);
    });
  }

  //
  // onStrikeChange - On strike change.
  //
  onStrikeChange(leg: TradeOptionLegs, index: number) {
    this.loadLegSymbol(leg);
  }

  //
  // onTypeChange - When we pick a call or a put
  //
  onTypeChange(leg: TradeOptionLegs, index: number) {
    this.loadLegSymbol(leg);
  }

  //
  // Duplicate leg
  //
  duplicateLeg(leg: TradeOptionLegs) {
    let newLeg = new TradeOptionLegs();
    newLeg.Symbol = leg.Symbol;
    newLeg.Expire = leg.Expire;
    newLeg.Type = leg.Type;
    newLeg.Strike = leg.Strike;
    newLeg.Side = leg.Side;
    newLeg.Qty = leg.Qty;
    newLeg.Strikes = leg.Strikes;
    this.tradeDetails.Legs.push(newLeg);

    // Load leg quotes
    this.loadLegSymbol(leg);
  }

  //
  // Remove leg
  //
  removeLeg(leg: TradeOptionLegs, index: number) {
    this.tradeDetails.Legs.splice(index, 1);
  }

  //
  // Load symbol for the leg
  //
  loadLegSymbol(leg: TradeOptionLegs) {
    // Ajax call to get option symbol based on these 3 params
    this.symbolService.getOptionSymbolFromParts(this.symbol.ShortName, leg.Expire, leg.Strike, leg.Type).subscribe(data => {
      leg.Symbol = data;
    });

  }  

}

/* End File */