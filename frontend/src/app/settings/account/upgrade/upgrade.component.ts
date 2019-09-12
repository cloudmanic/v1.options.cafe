//
// Date: 10/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { ActivatedRoute, Params } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { environment } from 'environments/environment';

const pageTitle: string = environment.title_prefix + "Settings Account Upgrade";

@Component({
	selector: 'app-settings-account-upgrade',
	templateUrl: './upgrade.component.html',
	styleUrls: []
})

export class UpgradeComponent implements OnInit {
	back: string = "";
	showCloseDownAccount: boolean = false;

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

			// Set the back
			if (params['back']) {
				this.back = params['back'];
			}

		});

	}

	//
	// Cancel account
	//
	cancelAccount() {
		this.showCloseDownAccount = true;
	}

	//
	// Cancel account
	//
	doCancelAccount() {
		this.showCloseDownAccount = false;
	}
}

/* End File */
