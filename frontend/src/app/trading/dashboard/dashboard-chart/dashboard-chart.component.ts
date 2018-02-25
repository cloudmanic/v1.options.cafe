//
// Date: 2/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { chart } from 'highcharts';
import * as Highcharts from 'highcharts';
import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { StateService } from '../../../providers/state/state.service';
import { QuotesService } from '../../../providers/http/quotes.service';

@Component({
  selector: 'app-trading-dashboard-chart',
  templateUrl: './dashboard-chart.component.html'
})

export class DashboardChartComponent implements OnInit 
{  
  @ViewChild('chartTarget') chartTarget: ElementRef;

  chart: Highcharts.ChartObject;
  
  symbol: string = "spy";

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
  // After View Init.
  //
  ngAfterViewInit() 
  {
    //
    // Setup the high charts options.
    //
    const options: Highcharts.Options = {
      chart: {
        type: 'bar'
      },
      title: {
        text: 'Fruit Consumption'
      },
      xAxis: {
        categories: ['Apples', 'Bananas', 'Oranges']
      },
      yAxis: {
        title: {
          text: 'Fruit eaten'
        }
      },
      series: [{
        name: 'Jane',
        data: [1, 0, 4]
      }, {
        name: 'John',
        data: [5, 7, 3]
      }]
    };
  
    // Load the chart
    this.chart = chart(this.chartTarget.nativeElement, options);
  }

  //
  // Add Series.
  //
  addSeries(){
    this.chart.addSeries({
      name:'Balram',
      data:[2,3,7]
    })    
  }  

  //
  // Update chart.
  //
  getChartData()
  {
    this.quoteService.getHistoricalQuote(this.symbol, new Date("2018-01-01"), new Date('2018-03-01'), 'daily').subscribe((res) => {
      console.log(res);
    });    
  }   

}

/* End File */