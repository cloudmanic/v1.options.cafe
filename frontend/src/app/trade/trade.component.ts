//
// Date: 4/2/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Symbol } from '../models/symbol';
import { Component, OnInit } from '@angular/core';
import { SymbolService } from '../providers/http/symbol.service';
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
  showTradeBuilder: boolean = false;
  tradeDetails: TradeDetails = new TradeDetails();

  //
  // Construct.
  //
  constructor(private websocketService: WebsocketService, private tradeService: TradeService, private optionsChainService: OptionsChainService, private symbolService: SymbolService) {

    // Default values
    this.tradeDetails.Class = "multileg";
    this.tradeDetails.Duration = "day";
    this.tradeDetails.OrderType = "market";

    // Subscribe to data updates from the quotes - Market Quotes
    this.websocketService.quotePushData.subscribe(data => {
      this.quotes[data.symbol] = data;
    });
  }

  //
  // OnInit.
  //
  ngOnInit() {
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
    // Successful Trade
    if(data.Action == "trade-success")
    {
      this.toggleTradeBuilder();
      return;
    }

    // Scroll to top of page
    window.scrollTo(0, 0);

    // Populate the form.
    this.tradeDetails = data.TradeDetails;

    // Set default if need be.
    if(this.tradeDetails.Class == undefined) 
    {
      this.tradeDetails.Class = "multileg";
      this.tradeDetails.Duration = "day";
      this.tradeDetails.OrderType = "market";
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
  setTradeClass(type: string) {
    this.tradeDetails.Class = type;
  }  

}

/* End File */