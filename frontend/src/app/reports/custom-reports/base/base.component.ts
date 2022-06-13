//
// Date: 3/9/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment-timezone';
import { Component, OnInit } from '@angular/core';
import { ReportType, Shared, TimeFrame } from 'app/reports/custom-reports/shared';
import { Router } from '@angular/router';
import { StateService } from 'app/providers/state/state.service';
import { Subject } from 'rxjs';

@Component({
	selector: 'app-base',
	template: ''
})

export class BaseComponent {
	groupBy: string = "month";

	// Date Stuff
	dateSelect: string = "ytd";
	startDate: Date = moment(moment().year() + "-01-01").toDate();
	endDate: Date = moment().toDate();
	startDateInput: String = moment(moment().year() + "-01-01").format('YYYY-MM-DD');
	endDateInput: String = moment().format('YYYY-MM-DD');
	dateTimeFrame: TimeFrame;
	dateTimeframes: TimeFrame[] = Shared.TimeFrames;

	// Setup options for reports type selector
	reportType: ReportType;
	reportTypes: ReportType[] = Shared.ReportTypes;

	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor.
	//
	constructor(public router: Router, public stateService: StateService) {
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
	// Set report type
	//
	setReportType(type: string) {
		for (let i = 0; i < this.reportTypes.length; i++) {
			if (this.reportTypes[i].Key == type) {
				this.reportType = this.reportTypes[i];
			}
		}
	}

	//
	// Change our custom report
	//
	doReportChange() {
		if (this.reportType.Query) {
			this.router.navigate([this.reportType.Route], { queryParams: { type: this.reportType.Key } });
		} else {
			this.router.navigate([this.reportType.Route]);
		}
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
	// We override this.
	//
	doBuildPage() { }
}

/* End File */
