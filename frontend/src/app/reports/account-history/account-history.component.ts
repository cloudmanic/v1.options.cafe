//
// Date: 9/12/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { BrokerEvent } from '../../models/broker-event';
import { Component, OnInit } from '@angular/core';
import { StateService } from '../../providers/state/state.service';
import { BrokerEventsService } from '../../providers/http/broker-events.service';
import { Title } from '@angular/platform-browser';
import { environment } from 'environments/environment';

const pageTitle: string = environment.title_prefix + "Reports Account History";

@Component({
	selector: 'app-account-history',
	templateUrl: './account-history.component.html',
	styleUrls: []
})

export class AccountHistoryComponent implements OnInit {
	page: number = 1;
	count: number = 0;
	limit: number = 0;
	noLimitCount: number = 0;
	events: BrokerEvent[];
	searchTerm: string = ""

	destory: Subject<boolean> = new Subject<boolean>();


	//
	// Construct.
	//
	constructor(private stateService: StateService, private brokerEventsService: BrokerEventsService, private titleService: Title) {
		this.page = this.stateService.GetAccountHistoryPage();
		this.events = this.stateService.GetAccountHistoryList();
	}

	//
	// NG Init
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Subscribe to changes in the selected broker.
		this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
			this.page = 1;
			this.getData();
		});

		// Load data for this screen
		this.getData();
	}

	//
	// OnDestroy
	//
	ngOnDestroy() {
		this.destory.next();
		this.destory.complete();
	}

	//
	// On paging click.
	//
	onPagingClick(page: number) {
		this.page = page;
		this.getData();
	}

	//
	// Get data.
	//
	getData() {
		// Make ajax call to get events.
		this.brokerEventsService.get(Number(this.stateService.GetStoredActiveAccountId()), 50, this.page, "date", "desc", "").subscribe((res) => {
			this.limit = res.Limit;
			this.noLimitCount = res.NoLimitCount;
			this.events = res.Data;
			this.count = res.Data.length;
			this.stateService.SetAccountHistoryList(res.Data);
			this.stateService.SetAccountHistoryPage(this.page);
		});
	}
}

/* End File */
