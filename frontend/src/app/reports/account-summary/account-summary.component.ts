//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment-timezone';
import * as Highcharts from 'highcharts/highstock';
import { Observable } from 'rxjs/Rx';
import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit } from '@angular/core';
import { SummaryYearly } from '../../models/reports';
import { StateService } from '../../providers/state/state.service';
import { ReportsService } from '../../providers/http/reports.service';
import { DropdownAction } from '../../shared/dropdown-select/dropdown-select.component';
import { environment } from 'environments/environment';
import { Title } from '@angular/platform-browser';

const pageTitle: string = environment.title_prefix + "Reports Account Summary";

@Component({
	selector: 'app-account-summary',
	templateUrl: './account-summary.component.html',
	styleUrls: []
})

export class AccountSummaryComponent implements OnInit {
	tradeGroupYears: number[];
	summaryByYear: SummaryYearly;
	summaryByYearSelected: number;
	summaryActions: DropdownAction[] = null;

	private destory: Subject<boolean> = new Subject<boolean>();

	// High charts config
	showFirstRun: boolean = false;

	Highcharts = Highcharts;

	groupBy: string = "month";

	startDate: Date = moment(moment().year() + "-01-01").toDate();

	endDate: Date = moment().toDate();

	chartType: string = "column";

	chartConstructor = 'chart';

	chartUpdateFlag: boolean = false;

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
	constructor(private stateService: StateService, private reportsService: ReportsService, private titleService: Title) {
		// Get data from site state.
		this.summaryByYear = this.stateService.GetReportsSummaryByYear();
		this.summaryByYearSelected = this.stateService.GetReportsSummaryByYearSelectedYear();
	}

	//
	// NG Init
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Subscribe to changes in the selected broker.
		this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
			this.summaryByYearSelected = new Date().getFullYear();
			this.stateService.SetReportsSummaryByYearSelectedYear(this.summaryByYearSelected);
			this.buildChart();
			this.getAccountSummary();
			this.getTradeGroupYears();
		});

		// Get data on page load.
		this.buildChart();
		this.getAccountSummary();
		this.getTradeGroupYears();
	}

	//
	// OnDestroy
	//
	ngOnDestroy() {
		this.destory.next();
		this.destory.complete();
	}

	//
	// Set chart type.
	//
	setChartType(type: string) {
		this.chartType = type;
		this.buildChart();
	}

	//
	// Get chart data
	//
	buildChart() {
		this.startDate = moment(this.summaryByYearSelected + "-01-01").toDate();
		this.endDate = moment(this.summaryByYearSelected + "-12-31").toDate();

		// Ajax call to get data
		this.reportsService.getProfitLoss(Number(this.stateService.GetStoredActiveAccountId()), this.startDate, this.endDate, this.groupBy, "asc", false).subscribe((res) => {

			// Show first run
			if (res.length > 0) {
				this.showFirstRun = false;
			} else {
				this.showFirstRun = true;
			}

			var data = [];

			for (var i = 0; i < res.length; i++) {
				let color = "#5cb85c";

				if (res[i].Profit < 0) {
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
	// Setup Summary actions.
	//
	setupSummaryActions() {
		let das = []
		this.summaryActions = []

		// Loop through add dates to the drop down
		for (let i = 0; i < this.tradeGroupYears.length; i++) {
			// First action
			let da = new DropdownAction();
			da.title = 'Year ' + this.tradeGroupYears[i];

			// Click on year
			da.click = (row: number[]) => {
				this.summaryByYearSelected = this.tradeGroupYears[i];
				this.getAccountSummary();
				this.stateService.SetReportsSummaryByYearSelectedYear(this.summaryByYearSelected);

				// Hack to get it to close;
				this.buildChart();
				this.setupSummaryActions();
			};

			das.push(da);
		}

		this.summaryActions = das;
	}

	//
	// Get Data = Trade Group Years
	//
	getTradeGroupYears() {
		// Make api call to get years
		this.reportsService.getTradeGroupYears(Number(this.stateService.GetStoredActiveAccountId())).subscribe((res) => {
			this.tradeGroupYears = res;
			this.setupSummaryActions();
		});
	}

	//
	// Get Data - Account Summary
	//
	getAccountSummary() {
		// Make api call to get account summary
		this.reportsService.getSummaryByYear(Number(this.stateService.GetStoredActiveAccountId()), this.summaryByYearSelected).subscribe((res) => {
			this.summaryByYear = res;
			this.stateService.SetReportsSummaryByYear(this.summaryByYear);
		});
	}

}

/* End File */
