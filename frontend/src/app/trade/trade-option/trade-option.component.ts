//
// Date: 11/13/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit, Input } from '@angular/core';
import { TradeDetails, TradeService, OrderPreview, TradeEvent } from '../../providers/http/trade.service';
import { SymbolService } from '../../providers/http/symbol.service';
import { StateService } from '../../providers/state/state.service';
import { Symbol } from '../../models/symbol';

@Component({
  selector: 'app-trade-option',
  templateUrl: './trade-option.component.html',
  styleUrls: []
})

export class TradeOptionComponent implements OnInit 
{
  @Input() quotes = {};
  @Input() showTradeBuilder: boolean;
  @Input() tradeDetails: TradeDetails;

  symbol: Symbol;
  option: Symbol;
  orderPreview: OrderPreview = null;
  typeAheadSymbol: Symbol;

  symbolType: string = "Call";
  symbolStrike: number = null;
  symbolStrikes: number[] = [];
  symbolExpire: Date = null;
  symbolExpirations: Date[] = [];

  //
  // Construct.
  //
  constructor(private tradeService: TradeService, private symbolService: SymbolService, private stateService: StateService) { }

  //
  // OnInit.
  //
  ngOnInit() 
  {
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
  // Reset trade
  //
  restTrade() 
  {
    // Set Defaults
    this.tradeDetails.Qty = 1;
    this.tradeDetails.Side = "buy_to_open";
    this.tradeDetails.Class = "option";
    this.tradeDetails.OrderType = "market";
    this.tradeDetails.Duration = "day";
    this.tradeDetails.Price = 0.00;

    // Reset symbol
    this.symbol = null;
    this.typeAheadSymbol = null;

    // Reset preview
    this.orderPreview = null;
  } 

  //
  // Submit Trade
  //
  submitTrade() 
  {
    // Set option symbol
    this.tradeDetails.OptionSymbol = this.option.ShortName;

    // Ajax call to submit trade.
    this.tradeService.submitTrade(this.tradeDetails, this.stateService.GetStoredActiveAccountId()).subscribe(

      // Success (as in no server errors)
      data => {
        if (data.Status == "ok")
        {
          this.restTrade();

          // Close trade window
          let event = new TradeEvent();
          event.Action = "trade-success";
          this.tradeService.tradeEvent.emit(event);

          // Show success notice
          this.stateService.SiteSuccess.emit("Order Submitted: Your order number is #" + data.Id); 
        } else
        {
          alert(data.Error);
        }
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
  // Preview Trade
  //
  previewTrade() 
  {
    // Set option symbol
    this.tradeDetails.OptionSymbol = this.option.ShortName;

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
    this.tradeDetails.Side = "buy_to_open";

    // Set the symbol
    this.tradeDetails.Symbol = this.symbol.ShortName;

    // Load expired dates
    this.loadExpireDates()

    // Make sure we are getting quotes from the websocket for this symbol.
    this.symbolService.addActiveSymbol(this.symbol.ShortName).subscribe();
  }

  //
  // Load expire dates
  //
  loadExpireDates() 
  {
    // Reset preview
    this.orderPreview = null;

    this.symbolExpirations = [];

    // Make API call to get option expire dates.
    this.tradeService.getOptionExpirations(this.tradeDetails.Symbol).subscribe(data => {
      
      // Set data
      this.symbolExpirations = data;

      // Set first expire date.
      if(this.symbolExpirations.length > 0)
      {
        this.symbolExpire = this.symbolExpirations[0];
      }

      // Load the strikes
      this.onExpireChange();

    });
  }

  //
  // On Strike change
  //
  onStrikeChange()
  {
    // Reset preview
    this.orderPreview = null;

    // Load option quotes
    this.loadOptionSymbol();
  }

  //
  // On type change
  //
  onTypeChange()
  {
    // Reset preview
    this.orderPreview = null;

    // Load option quotes
    this.loadOptionSymbol();
  }

  //
  // On expire change.
  //
  onExpireChange() 
  {
    // Reset preview
    this.orderPreview = null;

    // Api call to get the strikes for this chain.
    this.tradeService.getOptionStrikesBySymbolExpiration(this.tradeDetails.Symbol, this.symbolExpire).subscribe(data => {
      
      this.symbolStrikes = data;

      if(this.symbolStrikes.length > 0)
      {
        this.symbolStrike = this.symbolStrikes[0];
      }

      // Load option quotes
      this.loadOptionSymbol();
    });
  } 

  //
  // Load symbol for the option
  //
  loadOptionSymbol() 
  {
    // Ajax call to get option symbol based on these 3 params
    this.symbolService.getOptionSymbolFromParts(this.tradeDetails.Symbol, this.symbolExpire, this.symbolStrike, this.symbolType).subscribe(data => {
      this.option = data;
    });

  }         
}

/* End File */
