//
// Date: 4/2/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../providers/http/trade.service';



@Component({
  selector: 'app-trade',
  templateUrl: './trade.component.html'
})

export class TradeComponent implements OnInit 
{
  tradeDetails: TradeDetails = new TradeDetails();
  showTradeBuilder: boolean = true;

  //
  // Construct.
  //
  constructor(private tradeService: TradeService) 
  { 
    // Set Defaults (also used for development)
    this.tradeDetails.Symbol = "SPY";
    this.tradeDetails.Class = "multileg";
    this.tradeDetails.OrderType = "credit";
    this.tradeDetails.Duration = "gtc";
    this.tradeDetails.Price = 0.21;
    
    // Build legs
    this.tradeDetails.Legs = [];
    this.tradeDetails.Legs.push(new TradeOptionLegs().createNew("SPY180402P00262000", "sell_to_open", 10));
    this.tradeDetails.Legs.push(new TradeOptionLegs().createNew("SPY180402P00264000", "buy_to_open", 10));   
    

    console.log(this.tradeDetails); 
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

}

/* End File */