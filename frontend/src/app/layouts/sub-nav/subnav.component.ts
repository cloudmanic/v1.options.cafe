//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { ActivatedRoute } from '@angular/router';
import { TradeService, TradeEvent, TradeDetails } from '../../providers/http/trade.service';
import { Component, OnInit, ChangeDetectorRef } from '@angular/core';

@Component({
  selector: 'app-trading-sub-nav',
  templateUrl: './subnav.component.html'
})

export class SubnavComponent implements OnInit 
{
  routeData: any;
  action: string = '';
  section: string = 'trading';
  subSection: string = 'dashboard';

  //
  // Construct.
  //
  constructor(private route: ActivatedRoute, private tradeService: TradeService) { }

  //
  // OnInit
  //
  ngOnInit() 
  {
    this.routeData = this.route.data.subscribe(v => {
      this.action = v.action;      
      this.section = v.section;
      this.subSection = v.subSection;
    });    
  }

  //
  // OnDestory
  //
  ngOnDestroy() 
  {
    this.routeData.unsubscribe();
  }

  //
  // On Trade Click.
  //
  onTradeClick()
  {
    this.tradeService.PushEvent(new TradeEvent().createNew("toggle-trade-builder", new TradeDetails()));
  }

}

/* End File */