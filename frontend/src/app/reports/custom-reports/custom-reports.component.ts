//
// Date: 10/20/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment-timezone';
import * as Highcharts from 'highcharts/highstock';
import { ProfitLoss } from '../../models/reports';
import { Component, OnInit } from '@angular/core';
import { StateService } from '../../providers/state/state.service';
import { ReportsService } from '../../providers/http/reports.service';

@Component({
  selector: 'app-reports-custom-reports',
  templateUrl: './custom-reports.component.html',
  styleUrls: []
})

export class CustomReportsComponent implements OnInit 
{
  groupBy: string = "month";
  listData: ProfitLoss[] = [];
  Highcharts = Highcharts;

  chartConstructor = 'chart';

  chartUpdateFlag: boolean = false;

  // High charts config
  chartOptions = {

    chart: { type: 'column' },

    title: { text: '' },
    credits: { enabled: false },

    rangeSelector: { enabled: false },

    scrollbar: { enabled: false },

    navigator: { enabled: false },

    legend: { enabled: false },

    time: {
      getTimezoneOffset: function(timestamp) {
        // America/Los_Angeles
        // America/New_York
        let timezoneOffset = -moment.tz(timestamp, 'America/Los_Angeles').utcOffset();
        return timezoneOffset;
      }
    },

    tooltip: {
      formatter: function() {
        return "<b>" + this.points[0].series.name + ": </b><br />" + Highcharts.dateFormat('%b \'%y', this.points[0].x) + " : $" + Highcharts.numberFormat(this.points[0].y, 0, '.', ',');
      },

      shared: true
    },    

    yAxis: {
      title: {
        text: 'Profit & Loss'
      },

      labels: {
        formatter: function() {
          return '$' + Highcharts.numberFormat(this.axis.defaultLabelFormatter.call(this), 0, '.', ',');
        }
      }
    },    

    xAxis: {
      type: 'datetime',

      dateTimeLabelFormats: {
        month: '%b \'%y',
        year: '%Y',
        day: '%e. %b'
      },

      title: {
        text: 'Date'
      }
    },

    series: [{
      name: 'Profit & Loss',
      data: []
    }]
  };  

  //
  // Construct.
  //
  constructor(private stateService: StateService, private reportsService: ReportsService) 
  {

  }

  //
  // NG Init
  //
  ngOnInit() 
  {
    // Get data for page.
    this.buildChart();
    this.getProfitLoss();
  }

  //
  // Chart change`
  //
  chartChange()
  {
    this.buildChart();    
    this.getProfitLoss();
  }

  //
  // Get chart data
  //
  buildChart()
  {
    this.reportsService.getProfitLoss(Number(this.stateService.GetStoredActiveAccountId()), "2017-01-01", "2018-12-31", this.groupBy, "asc").subscribe((res) => {

      var data = [];

      for (var i = 0; i < res.length; i++) 
      {
        let color = "#5cb85c";

        if (res[i].Profit < 0) 
        {
          color = "#ce4260";
        }

        data.push({ x: res[i].Date, y: res[i].Profit, color: color });
      }

      // Rebuilt the chart
      this.chartOptions.series[0].data = data;
      this.chartOptions.series[0].name = "Profit & Loss";
      this.chartUpdateFlag = true;

    });
  }

  //
  // Get Data = Profit Loss
  //
  getProfitLoss() 
  {
    this.reportsService.getProfitLoss(Number(this.stateService.GetStoredActiveAccountId()), "2017-01-01", "2018-12-31", this.groupBy, "desc").subscribe((res) => {
      this.listData = res;
    });
  }

  //
  // Get profit total
  //
  getProfitTotal(rows: ProfitLoss[]): number 
  {
    let total = 0;

    for(let i = 0; i < rows.length; i++)
    {
      total += rows[i].Profit;
    }

    return total;
  }

  //
  // Get trade total
  //
  getTradeTotal(rows: ProfitLoss[]): number 
  {
    let total = 0;

    for (let i = 0; i < rows.length; i++) 
    {
      total += rows[i].TradeCount;
    }

    return total;
  }

  //
  // Get trade total
  //
  getCommissionsTotal(rows: ProfitLoss[]): number 
  {
    let total = 0;

    for (let i = 0; i < rows.length; i++) 
    {
      total += rows[i].Commissions;
    }

    return total;
  }

}

/* End File */