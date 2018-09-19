//
// Date: 9/19/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { HttpErrorResponse } from '@angular/common/http';
import { Symbol } from '../../models/symbol';
import { Component, OnInit, Input } from '@angular/core';
import { SymbolService } from '../../providers/http/symbol.service';
import { StateService } from '../../providers/state/state.service';
//import { OptionsChainService } from '../../providers/http/options-chain.service';
import { TradeService, TradeEvent, TradeDetails, OrderPreview } from '../../providers/http/trade.service';

@Component({
  selector: 'app-trade-equity',
  templateUrl: './trade-equity.component.html',
  styleUrls: []
})

export class TradeEquityComponent implements OnInit 
{
  @Input() quotes = {};
  @Input() showTradeBuilder: boolean;
  @Input() tradeDetails: TradeDetails;

  symbol: Symbol;
  orderPreview: OrderPreview;
  typeAheadSymbol: Symbol;
  symbolExpirations: Date[] = [];

  //
  // Construct.
  //
  constructor(private tradeService: TradeService, private symbolService: SymbolService, private stateService: StateService) { }

  //
  // OnInit.
  //
  ngOnInit() {

    // Get and set symbols
    if (this.tradeDetails.Symbol) 
    {
      this.symbolService.getSymbol(this.tradeDetails.Symbol).subscribe(data => {
        this.symbol = data;
        this.typeAheadSymbol = data;
      });
    }

  }

  //
  // Preview Trade
  //
  previewTrade() 
  {
    // Ajax call to preview trade.
    this.tradeService.previewTrade(this.tradeDetails, this.stateService.GetStoredActiveAccountId()).subscribe(

      // Success (as in no server errors)
      data => {
        this.orderPreview = data;
      },

      // Error
      (err: HttpErrorResponse) => {

        if (err.error instanceof Error) 
        {
          alert(err.error.message);
        } else 
        {
          alert(err.error.error);
        }

      }

    );
  }  

  //
  // Reset trade
  //
  restTrade() 
  {
    // Set Defaults
    this.tradeDetails.Qty = 1;
    this.tradeDetails.Side = "buy";
    this.tradeDetails.Class = "equity";
    this.tradeDetails.OrderType = "market";
    this.tradeDetails.Duration = "day";
    this.tradeDetails.Price = 0.00;

    // Reset symbol
    this.symbol = null;
    this.typeAheadSymbol = null;
  }  

  //
  // onSearchTypeAheadClick() 
  //
  onSearchTypeAheadClick(symbol: Symbol) 
  {
    if(typeof symbol == "undefined") 
    {
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

    // Make sure we are getting quotes from the websocket for this symbol.
    this.symbolService.addActiveSymbol(this.symbol.ShortName).subscribe();
  }

}

/* End File */
