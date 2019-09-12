//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Router } from '@angular/router';
import { BrokerService } from '../../providers/http/broker.service';
import { environment } from '../../../environments/environment';
import { Title } from '@angular/platform-browser';

const pageTitle: string = environment.title_prefix + "Broker Select";

@Component({
	selector: 'broker-select',
	templateUrl: './home.component.html'
})

export class AuthBrokerSelectComponent implements OnInit {

	broker: string = "tradier"

	constructor(private router: Router, private brokerService: BrokerService, private titleService: Title) { }

	//
	// On Init
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// No user id redirect
		if (!localStorage.getItem('user_id').length) {
			this.router.navigate(['/login']);
		}

		// No access token redirect
		if (!localStorage.getItem('access_token').length) {
			this.router.navigate(['/login']);
		}

	}

	//
	// Login submit.
	//
	onSubmit(form: NgForm) {

		// Ajax call to add broker.
		this.brokerService.create("Tradier", "Tradier Account").subscribe((res) => {

			// Switch based on broker selected - Redirect to login to broker and get access token.
			switch (form.value["field-broker"]) {
				case 'tradier':
					window.location.href = environment.app_server + '/tradier/authorize?user=' + localStorage.getItem('user_id') + '&broker_id=' + res.Id;
					break;
			}

		});

	}

}

/* End File */
