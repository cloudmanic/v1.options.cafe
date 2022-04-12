import { Component, OnInit } from '@angular/core';
import { environment } from 'environments/environment';
import { Title } from '@angular/platform-browser';

const pageTitle: string = environment.title_prefix + "Create Backtest";

@Component({
	selector: 'app-backtest-create',
	templateUrl: './create.component.html'
})

export class BacktestCreateComponent implements OnInit {
	//
	// Constructor
	//
	constructor(private titleService: Title) { }

	//
	// ngOninit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);
	}
}

/* End File */
