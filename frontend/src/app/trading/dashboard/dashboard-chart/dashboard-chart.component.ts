//
// Date: 2/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Rx';
import * as Highcharts from 'highcharts/highstock';
import 'rxjs/add/operator/takeUntil';
import * as moment from 'moment-timezone';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { Symbol } from '../../../models/symbol';
import { StateService } from '../../../providers/state/state.service';
import { QuotesService } from '../../../providers/http/quotes.service';
import { WebsocketService } from '../../../providers/http/websocket.service';

@Component({
  selector: 'app-trading-dashboard-chart',
  templateUrl: './dashboard-chart.component.html'
})

export class DashboardChartComponent implements OnInit 
{  
  quotes = {}

  symbol: Symbol = new Symbol().New(1, "SPDR S&P 500 ETF Trust", "SPY", "Equity", null, null, null, null);
  interval: string = "daily";
  rangeSelect: string;
  destory: Subject<boolean> = new Subject<boolean>();

  Highcharts = Highcharts;

  chartConstructor = 'stockChart';
  
  chartUpdateFlag: boolean = false;
  
  // High charts config
  chartOptions = {
    title: { text: '' },
    credits: { enabled: false },

    rangeSelector: { enabled: false },
    
    scrollbar: { enabled: false },

    navigator: { enabled: false },

    legend: { enabled: false },

    time: {
      getTimezoneOffset: function (timestamp) {
        // America/Los_Angeles
        // America/New_York
        let timezoneOffset = -moment.tz(timestamp, 'America/Los_Angeles').utcOffset();
        return timezoneOffset;
      }
    },

    yAxis: {
      crosshair: true,
      startOnTick: false,
      endOnTick: false,
      minPadding: 0.1,
      maxPadding: 0.1          
    },  

    xAxis : {
      type: 'datetime',
      crosshair: true,
      minRange: 3600 * 1000,
      dateTimeLabelFormats: {
        second: '%I:%M:%S %p',
        minute: '%I:%M %p',
        hour: '%I:%M %p',
        day: '%m/%e %I:%M %p',
        week: '%m/%e %I:%M %p',
        month: '%m/%Y',
        year: '%Y'
      },          
    },              

    series : [{
      name : 'SPY',
      type: 'candlestick',
      data: [],
      turboThreshold: 0,
      tooltip: { valueDecimals: 2 },
      dataGrouping: { enabled: false }
    }]
  };

  //
  // Constructor....
  //
  constructor(
    private stateService: StateService, 
    private quotesService: QuotesService, 
    private websocketService: WebsocketService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    // Get data from cache.
    this.quotes = this.stateService.GetQuotes();     
    this.symbol = this.stateService.GetDashboardChartSymbol();
    this.rangeSelect = this.stateService.GetDashboardChartRangeSelect();    
    this.chartOptions.series[0].name = this.symbol.ShortName;
    this.chartOptions.series[0].data = this.stateService.GetDashboardChartData();

    // Subscribe to data updates from the quotes - Market Quotes
    this.websocketService.quotePushData.takeUntil(this.destory).subscribe(data => {
      this.quotes[data.symbol] = data;
    });     

    // Reload the chart every 1min after a 1 min delay to start
    Observable.timer((1000 * 60), (1000 * 60)).takeUntil(this.destory).subscribe(x => { this.getChartData(); });

    // Load data for the page.
    this.getChartData();
  }

  //
  // OnDestroy
  //
  ngOnDestroy()
  {
    this.destory.next();
    this.destory.complete();
  }  

  //
  // onSearchTypeAheadClick() 
  //
  onSearchTypeAheadClick(symbol: Symbol) {

    if(typeof symbol == "undefined")
    {
      return;
    }

    this.symbol = symbol;
    this.getChartData();
    this.stateService.SetDashboardChartSymbol(symbol);    
  }  

  //
  // Update chart.
  //
  getChartData()
  {
    // Make api call to get historical data.
    this.quotesService.getHistoricalQuote(this.symbol.ShortName, this.getStartDate(), new Date(), this.interval).subscribe((res) => {
      var data = [];
      
      for(var i = 0; i < res.length; i++)
      {
        data.push({
          x: res[i].Date,
          open: res[i].Open,
          high: res[i].High,
          low: res[i].Low,
          close: res[i].Close,
          name: (res[i].Date.getMonth() + 1) + "/" + res[i].Date.getDate() +  "/" + res[i].Date.getFullYear(),
          color: (((res[i].Close - res[i].Open) > 0) ? '#5cb85c' : '#ce4260')
        });
      }

      // Rebuilt the chart
      this.chartOptions.series[0].data = data;
      this.chartOptions.series[0].name = this.symbol.ShortName;
      this.chartUpdateFlag = true;

      // Store cache
      this.stateService.SetDashboardChartData(data);
      this.stateService.SetDashboardChartRangeSelect(this.rangeSelect);
    });    
  }

  //
  // Get start date.
  //
  getStartDate() : Date 
  {
    // See if today is a weekend
    let start = new Date();
    let dayCount = moment().isoWeekday();

    switch(dayCount)
    {
      case 6:
        start = moment().subtract(1, 'day').toDate();
      break;

      case 7:
        start = moment().subtract(2, 'days').toDate();
      break;    
    }

    let parts = this.rangeSelect.split("-");
    this.interval = parts[1];

    switch(parts[0])
    {
      case 'days':
        let numberOfDaysToSubtract: any = parts[2];
        start.setDate(start.getDate() - numberOfDaysToSubtract);        
      break;
    }

    return start;
  }   

  //
  // Change date range.
  //
  onRangeSelect(event) 
  {
    this.getChartData();
  }

}

/* End File */