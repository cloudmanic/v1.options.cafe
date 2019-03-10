//
// Date: 10/20/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { StateService } from 'app/providers/state/state.service';
import { Router, ActivatedRoute, Params } from '@angular/router';

@Component({
	selector: 'app-reports-custom-reports',
	templateUrl: './custom-reports.component.html',
	styleUrls: []
})

export class CustomReportsComponent {
	dataType: string = "";

	//
	// Construct.
	//
	constructor(private router: Router, private stateService: StateService, private activatedRoute: ActivatedRoute) {
		// subscribe to router event
		this.activatedRoute.queryParams.subscribe((params: Params) => {
			// See what type profit / loss we have
			if (params['type']) {
				this.dataType = params['type'];
			}

			// Redirect to profit / loss page.
			if (this.dataType.length) {
				this.router.navigate(['/reports/custom/profit-loss'], { queryParams: { type: this.dataType } });
			} else {
				this.router.navigate(['/reports/custom/profit-loss']);
			}
		});
	}
}

/* End File */
