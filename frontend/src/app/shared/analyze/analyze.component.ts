//
// Date: 11/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import * as moment from 'moment-timezone';
import * as Highcharts from 'highcharts/highstock';
import { Observable } from 'rxjs/Rx';
import { Subject } from 'rxjs/Subject';
import { faTimes } from '@fortawesome/free-solid-svg-icons';
import { AnalyzeService, AnalyzeTrade } from '../../providers/http/analyze.service'
import { Component,  OnInit, Input } from '@angular/core';
import { AnalyzeResult, AnalyzeLeg } from '../../models/analyze-result';

@Component({
  selector: 'app-shared-analyze',
  templateUrl: './analyze.component.html',
  styleUrls: []
})

export class AnalyzeComponent implements OnInit 
{
  private destory: Subject<boolean> = new Subject<boolean>(); 

  showChart: boolean = false;

  closeIcon = faTimes; 

  // High charts config
  Highcharts = Highcharts;

  chartUpdateFlag: boolean = false;

  chartOptions = {

    chart: { 
      type: 'line', 
      zoomType: 'x',
      panning: true,
      panKey: 'shift'    
    },

    title: {
      text: ''
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
        dashStyle: 'ShortDash',        
        color: '#E0A300',
        width: 2,
        value: 273.73
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
    // Get signals to open chart
    this.analyzeService.dialog.takeUntil(this.destory).subscribe((trade: AnalyzeTrade) => {
      this.getResults(trade);
      this.showChart = true;
    });
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
  // Close Dialog
  //
  closeDialog() 
  {
    this.showChart = false;
  }

  //
  // Get analyze results
  //
  getResults(trade: AnalyzeTrade)
  {
    // Make Ajax call to get chart data
    this.analyzeService.getOptionsUnderlyingPriceResult(trade.OpenCost, trade.Legs).subscribe((res) => {

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