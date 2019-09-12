//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { environment } from 'environments/environment';
import { Title } from '@angular/platform-browser';

const pageTitle: string = environment.title_prefix + "Account Upgrade Credit Card";

@Component({
	selector: 'app-settings-account-upgrade-credit-card',
	templateUrl: './credit-card.component.html',
	styleUrls: []
})
export class CreditCardComponent implements OnInit {
	plan: string = "";
	today: Date = new Date();

	//
	// Construct.
	//
	constructor(private activatedRoute: ActivatedRoute, private titleService: Title) { }

	//
	// OnInit...
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// subscribe to router event
		this.activatedRoute.queryParams.subscribe((params: Params) => {

			// Set the plan
			if (params['plan']) {
				this.plan = params['plan'];
			}

		});

	}

}

/* End File */
