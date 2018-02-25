//
// Date: 2/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import * as Highcharts from 'highcharts/highstock';
import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { StateService } from '../../../providers/state/state.service';
import { QuotesService } from '../../../providers/http/quotes.service';

@Component({
  selector: 'app-trading-dashboard-chart',
  templateUrl: './dashboard-chart.component.html'
})

export class DashboardChartComponent implements OnInit 
{  
  symbol: string = "spy";

  Highcharts = Highcharts;

  chartConstructor = 'stockChart';
  
  chartUpdateFlag: boolean = false;
  
  // High charts config
  chartOptions = {
    title: { text: '' },
    credits: { enabled: false },

    rangeSelector: { enabled: false },

    yAxis: {
      startOnTick: false,
      endOnTick: false,
      minPadding: 0.1,
      maxPadding: 0.1          
    },  

    xAxis : {
      type: 'datetime',
      minRange: 3600 * 1000 // one hour
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
  constructor(private stateService: StateService, private quoteService: QuotesService) { }

  //
  // OnInit....
  //
  ngOnInit() 
  {
    this.getChartData()
  }

  //
  // Update chart.
  //
  getChartData()
  {
    // Make api call to get historical data.
    this.quoteService.getHistoricalQuote(this.symbol, new Date("2018-01-01"), new Date('2018-03-01'), 'daily').subscribe((res) => {
      var data = [];
      
      for(var i = 0; i < res.length; i++)
      {
        data.push({
          x: res[i].Date,
          open: res[i].Open,
          high: res[i].High,
          low: res[i].Low,
          close: res[i].Close,
          name: (res[i].Date.getMonth() + 1) + "/" + res[i].Date.getDay() +  "/" + res[i].Date.getFullYear()
          //color: '#00FF00'
        });
      }

      // Rebuilt the chart
      this.chartOptions.series[0].data = data;
      this.chartUpdateFlag = true;
    });    
  }   

}

/* End File */