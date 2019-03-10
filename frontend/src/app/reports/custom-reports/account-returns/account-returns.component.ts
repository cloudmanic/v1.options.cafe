//
// Date: 3/9/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment-timezone';
import * as Highcharts from 'highcharts/highstock';
import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { Angular2Csv } from 'angular2-csv/Angular2-csv';
import { ProfitLoss, AccountReturn } from '../../../models/reports';
import { Component, OnInit } from '@angular/core';
import { StateService } from '../../../providers/state/state.service';
import { ReportsService } from '../../../providers/http/reports.service';
import { Router } from '@angular/router';

@Component({
	selector: 'app-reports-custom-reports-account-returns',
	templateUrl: './account-returns.component.html'
})
export class AccountReturnsComponent implements OnInit {
	dataType: string = "account-returns";
	showFirstRun: boolean = false;
	dateSelect: string = "ytd";
	chartType: string = "column";
	groupBy: string = "month";
	startDate: Date = moment(moment().year() + "-01-01").toDate();
	endDate: Date = moment().toDate();
	startDateInput: Date = moment(moment().year() + "-01-01").format('YYYY-MM-DD');
	endDateInput: Date = moment().format('YYYY-MM-DD');

	arData: AccountReturn[] = [];

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
				return "<b>" + this.points[0].series.name + ": </b><br />" + Highcharts.dateFormat('%b \'%y', this.points[0].x) + " : " + Highcharts.numberFormat(this.points[0].y, 0, '.', ',') + "%";
			},

			shared: true
		},

		yAxis: {
			title: {
				text: '% Gain / Loss'
			},

			labels: {
				formatter: function() {
					return Highcharts.numberFormat(this.axis.defaultLabelFormatter.call(this), 0, '.', ',') + '%';
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
			name: 'Account Returns',
			data: []
		}]
	};

	//
	// Construct.
	//
	constructor(private router: Router, private stateService: StateService, private reportsService: ReportsService) {
		// Subscribe to changes in the selected broker.
		this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
			this.doBuildPage();
		});
	}

	//
	// OnInit....
	//
	ngOnInit() {
		// Load page.
		this.doBuildPage();
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
		this.doBuildPage();
	}

	//
	// Data type change.
	//
	dataTypeChange() {
		switch (this.dataType) {
			case "profit-loss":
				this.router.navigate(['/reports/custom/profit-loss']);
				break;

			case "profit-loss-cumulative":
				this.router.navigate(['/reports/custom'], { queryParams: { type: "profit-loss-cumulative" } });
				break;

			case "account-returns":
				this.doBuildPage();
				break;
		}
	}

	//
	// Do account returns graph
	//
	doBuildPage() {
		// AJAX call to get data`
		this.reportsService.getAccountReturns(Number(this.stateService.GetStoredActiveAccountId()), this.startDate, this.endDate).subscribe((res) => {
			// Show first run
			if (res.length > 0) {
				this.showFirstRun = false;
			} else {
				this.showFirstRun = true;
				return;
			}

			// Set data
			this.arData = res;

			// Build chart
			var data = [];

			for (var i = 0; i < res.length; i++) {
				let color = "#5cb85c";

				if (res[i].Percent < 0) {
					color = "#ce4260";
				}

				data.push({ x: res[i].Date, y: (res[i].Percent * 100), color: color });
			}

			// Rebuilt the chart
			this.chartOptions.chart.type = this.chartType;
			this.chartOptions.series[0].data = data;
			this.chartOptions.series[0].name = "Account Returns";
			this.chartUpdateFlag = true;
		});
	}

	//
	// Deal with date change
	//
	onDateChange() {
		// Set start and stop dates based on predefined selector
		switch (this.dateSelect) {
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

		// Rebuilt the page.
		this.doBuildPage();
	}

	//
	// Export CSV
	//
	doExportCSV() {
		let data = [];

		// Build data
		for (let i = 0; i < this.arData.length; i++) {
			let row = this.arData[i];

			data.push({
				Date: moment(row.Date).format('YYYY-MM-DD'),
				Percent: row.Percent,
				TotalCash: row.TotalCash,
				AccountValue: row.AccountValue,
				PricePer: row.PricePer,
				Units: row.Units
			});
		}

		let options = {
			fieldSeparator: ',',
			quoteStrings: '"',
			decimalseparator: '.',
			headers: ['Date', 'Percent', 'TotalCash', 'AccountValue', 'PricePer', 'Units'],
			showTitle: false,
			useBom: true,
			removeNewLines: false,
			keys: []
		};

		new Angular2Csv(data, 'options-cafe-account-returns', options);
	}
}

/* End File */
