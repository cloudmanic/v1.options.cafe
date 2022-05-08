import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Backtest } from 'app/models/backtest';
import { BacktestService } from 'app/providers/http/backtest.service';
import { environment } from 'environments/environment';

const pageTitle: string = environment.title_prefix + "Backtest View";

@Component({
	selector: 'app-backtest-view',
	templateUrl: './view.component.html'
})

export class BacktestViewComponent implements OnInit {
	backtest: Backtest = new Backtest();
	backtestId: number = 0;

	//
	// Constructor
	//
	constructor(private titleService: Title, private backtestService: BacktestService, private route: ActivatedRoute,) { }

	//
	// ngOninit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Backtest Id
		this.backtestId = this.route.snapshot.params['id'];

		// Load page data
		this.getData();
	}

	//
	// Get list of backtests`
	//
	getData() {
		this.backtestService.getById(this.backtestId).subscribe(data => {
			this.backtest = data;

			console.log(this.backtest)
		});
	}

}

/* End File */
