import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { Backtest } from 'app/models/backtest';
import { BacktestService } from 'app/providers/http/backtest.service';
import { environment } from 'environments/environment';

const pageTitle: string = environment.title_prefix + "Backtest View";

@Component({
	selector: 'app-backtest-view',
	templateUrl: './view.component.html'
})

export class BacktestViewComponent implements OnInit {
	backtests: Backtest[] = [];

	//
	// Constructor
	//
	constructor(private titleService: Title, private backtestService: BacktestService) { }

	//
	// ngOninit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Load page data
		this.getData();
	}

	//
	// Get list of backtests`
	//
	getData() {
		this.backtestService.get().subscribe(data => {
			this.backtests = data;
		});
	}

}

/* End File */
