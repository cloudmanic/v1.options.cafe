//
// Date: 10/20/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment-timezone';
import * as Highcharts from 'highcharts/highstock';
import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { Angular2Csv } from 'angular2-csv/Angular2-csv';
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
  cumulative: boolean = false;
  dataType: string = "profit-loss";
  showFirstRun: boolean = false;
  dateSelect: string = "1-year";
  chartType: string = "column";
  groupBy: string = "month";
  startDate: Date = moment(moment().year() + "-01-01").toDate();
  endDate: Date = moment().toDate();
  startDateInput: Date = moment(moment().year() + "-01-01").format('YYYY-MM-DD');
  endDateInput: Date = moment().format('YYYY-MM-DD');

  listData: ProfitLoss[] = [];
  Highcharts = Highcharts;

  chartConstructor = 'chart';

  chartUpdateFlag: boolean = false;

  destory: Subject<boolean> = new Subject<boolean>();

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
        text: ''
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
  // Set chart type.
  //
  setChartType(type: string)
  {
    this.chartType = type;
    this.buildChart();
  }

  //
  // NG Init
  //
  ngOnInit() 
  {
    // Subscribe to changes in the selected broker.
    this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
      this.buildChart();
      this.getProfitLoss();
    });

    // Get data for page.
    this.buildChart();
    this.getProfitLoss();
  }

  //
  // OnDestroy
  //
  ngOnDestroy() {
    this.destory.next();
    this.destory.complete();
  }


  //
  // Data type change.
  //
  dataTypeChange()
  {
    switch(this.dataType)
    {
      case "profit-loss":
        this.cumulative = false;
      break;

      case "profit-loss-cumulative":
        this.cumulative = true;
      break;      
    }

    this.chartChange();
  }

  //
  // Chart change
  //
  chartChange()
  {
    this.buildChart();    
    this.getProfitLoss();
  }

  //
  // Deal with date change
  //
  dateChange()
  {
    // Set start and stop dates based on predefined selector
    switch(this.dateSelect)
    {
      case "1-year":
        this.startDate = moment().subtract(1, 'year').toDate();
        this.endDate = moment().toDate();
      break;

      case "2-year":
        this.startDate = moment().subtract(2, 'year').toDate();
        this.endDate = moment().toDate();
      break;

      case "3-year":
        this.startDate = moment().subtract(3, 'year').toDate();
        this.endDate = moment().toDate();
      break;

      case "4-year":
        this.startDate = moment().subtract(4, 'year').toDate();
        this.endDate = moment().toDate();
      break;

      case "5-year":
        this.startDate = moment().subtract(5, 'year').toDate();
        this.endDate = moment().toDate();
      break;

      case "10-year":
        this.startDate = moment().subtract(10, 'year').toDate();
        this.endDate = moment().toDate();
      break;

      case "ytd":
        this.startDate = moment(moment().year() + "-01-01").toDate();
        this.endDate = moment().toDate();
      break;

      case "custom":
        this.startDate = moment(this.startDateInput).toDate();
        this.endDate = moment(this.endDateInput).toDate();
      break;
    }

    this.buildChart();
    this.getProfitLoss();    
  }

  //
  // Get chart data
  //
  buildChart()
  {
    this.reportsService.getProfitLoss(Number(this.stateService.GetStoredActiveAccountId()), this.startDate, this.endDate, this.groupBy, "asc", this.cumulative).subscribe((res) => {

      // Show first run
      if(res.length > 0) 
      {
        this.showFirstRun = false;
      } else 
      {
        this.showFirstRun = true;
      }

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
      this.chartOptions.chart.type = this.chartType;
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
    let sort = "desc";

    if(this.cumulative)
    {
      sort = "asc";
    }

    this.reportsService.getProfitLoss(Number(this.stateService.GetStoredActiveAccountId()), this.startDate, this.endDate, this.groupBy, sort, this.cumulative).subscribe((res) => {
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

  //
  // Export CSV
  //
  exportCSV()
  {
    let data = [];

    // Build data 
    for(let i = 0; i < this.listData.length; i++)
    {
      let row = this.listData[i];

      data.push({
        Date: moment(row.Date).format('YYYY-MM-DD'),
        Profit: row.Profit,
        TradeCount: row.TradeCount,
        Commissions: row.Commissions,
        ProfitPerTrade: row.ProfitPerTrade,
        WinRatio: row.WinRatio,
        LossCount: row.LossCount,
        WinCount: row.WinCount        
      });
    }

    let options = {
      fieldSeparator: ',',
      quoteStrings: '"',
      decimalseparator: '.',
      headers: ['Date', 'Profit', 'TradeCount', 'Commissions', 'ProfitPerTrade', 'WinRatio', 'LossCount', 'WinCount'],
      showTitle: false,
      useBom: true,
      removeNewLines: false,
      keys: []
    };

    new Angular2Csv(data, 'options-cafe-profit-loss', options);
  }
}

/* End File */