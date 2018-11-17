//
// Date: 11/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment-timezone';
import * as Highcharts from 'highcharts/highstock';
import { AnalyzeService } from '../../providers/http/analyze.service'
import { Component, OnInit } from '@angular/core';
import { AnalyzeResult, AnalyzeLeg } from '../../models/analyze-result';

@Component({
  selector: 'app-shared-analyze',
  templateUrl: './analyze.component.html',
  styleUrls: []
})

export class AnalyzeComponent implements OnInit 
{
  // High charts config
  Highcharts = Highcharts;

  chartUpdateFlag: boolean = false;

  chartOptions = {

    chart: { 
      type: 'line', 
      zoomType: 'x',
      panning: true,
      panKey: 'shift',
      backgroundColor: {
        linearGradient: { x1: 0, y1: 0, x2: 1, y2: 1 },
        stops: [
            [0, 'rgb(255, 255, 255)'],
            [1, 'rgb(240, 240, 255)']
        ]
      },
      borderWidth: 2,
      plotBackgroundColor: 'rgba(255, 255, 255, .9)',
      plotShadow: true,
      plotBorderWidth: 1     
    },

    title: {
      text: 'Profit & Loss of Trade at Expiration'
    },

    subtitle: {
      text: ''
    },

    credits: { enabled: false },

    rangeSelector: { enabled: false },

    scrollbar: { enabled: false },

    navigator: { enabled: false },

    legend: { enabled: false },

    yAxis: {
      title: {
        text: 'Profit & Loss'
      },

      labels: {
        formatter: function() {
          return '$' + Highcharts.numberFormat(this.axis.defaultLabelFormatter.call(this), 0, '.', ',');
        }
      },

      plotLines: [{
        dashStyle: 'ShortDash',
        color: '#000000',
        width: 2,
        value: 0
      }]
           
    },    

    xAxis: {
      labels: {
        format: '${value}'
      },
      
      crosshair: true,

      title: {
        text: 'Underlying Price'
      },

      plotLines: [{
        color: '#FF0000',
        width: 2,
        value: 5.5
      }]    

    },
    
    tooltip: {
      headerFormat: '<b>Profit & Loss:</b> ${point.y:.2f}<br /><b>Underlying Price:</b> ${point.x}',
      pointFormat: '',
    },

    series: [{
      name: 'Profit & Loss',
      data: []
    }]
  };

  //
  // Construct.
  //
  constructor(private analyzeService: AnalyzeService) { }

  //
  // NgInit
  //
  ngOnInit() 
  {
    this.getResults();
  }

  //
  // Get analyze results
  //
  getResults()
  {
    let leg1 = new AnalyzeLeg();
    leg1.Qty = 1;
    leg1.SymbolStr = "SPY181221C00250000";

    let leg2 = new AnalyzeLeg();
    leg2.Qty = -2;
    leg2.SymbolStr = "SPY181221C00260000";

    let leg3 = new AnalyzeLeg();
    leg3.Qty = 1;
    leg3.SymbolStr = "SPY181221C00270000";    
    
    let legs: AnalyzeLeg[] = [ leg1, leg2, leg3 ];


    this.analyzeService.getOptionsUnderlyingPriceResult(157.00, 273.73, legs).subscribe((res) => {

      var data = [];

      for (var i = 0; i < res.length; i++) 
      {
        let color = "#5cb85c";

        if (res[i].Profit < 0) 
        {
          color = "#ce4260";
        }

        data.push({ x: res[i].UnderlyingPrice, y: res[i].Profit, color: color });
      }

      // Rebuilt the chart
      this.chartOptions.series[0].data = data;
      this.chartOptions.series[0].name = "Profit & Loss";
      this.chartUpdateFlag = true;

    });
  } 

}

/* End File */