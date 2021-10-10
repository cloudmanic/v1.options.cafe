//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Rx';
import { Subject } from 'rxjs/Subject';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Screener } from '../../models/screener';
import { ScreenerResult } from '../../models/screener-result';
import { StateService } from '../../providers/state/state.service';
import { ScreenerService } from '../../providers/http/screener.service';
import { faListAlt, faTh, faCaretRight, faCaretDown } from '@fortawesome/free-solid-svg-icons';
import { TradeService, TradeEvent, TradeDetails, TradeOptionLegs } from '../../providers/http/trade.service';
import { Settings } from 'app/models/settings';
import { SettingsService } from 'app/providers/http/settings.service';
import { environment } from 'environments/environment';
import { Title } from '@angular/platform-browser';

const pageTitle: string = environment.title_prefix + "Screener";

@Component({
	selector: 'app-screener',
	templateUrl: './home.component.html'
})

export class ScreenerComponent implements OnInit {
	screeners: Screener[] = [];
	destory: Subject<boolean> = new Subject<boolean>();
	listGrid = faTh;
	listIcon = faListAlt;
	faCaretDown = faCaretDown;
	faCaretRight = faCaretRight;
	settings: Settings = new Settings();

	//
	// Constructor....
	//
	constructor(private stateService: StateService, private screenerService: ScreenerService, private tradeService: TradeService, private router: Router, private settingsService: SettingsService, private titleService: Title) {
		// Load settings
		this.loadSettingsData();
	}

	//
	// OnInit....
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Default start timer
		let startTimer: number = (1000 * 10);

		// Get Data from cache
		this.screeners = this.stateService.GetScreens();

		// Load page data.
		if (!this.screeners) {
			startTimer = (1000 * 60);
			this.getScreeners();
		}

		// Reload the data every 1min after a 1 min delay to start
		Observable.timer(startTimer, (1000 * 60)).takeUntil(this.destory).subscribe(x => { this.getScreeners(); });
	}

	//
	// OnDestroy
	//
	ngOnDestroy() {
		this.destory.next();
		this.destory.complete();
	}

	//
	// Load settings data.
	//
	loadSettingsData() {
		this.settingsService.get().subscribe((res) => {
			this.settings = res;
			this.stateService.SetSettings(res);
		});
	}

	//
	// Sort click
	//
	sortClick(screen: Screener, col: string) {
		if (screen.ListSort == col) {
			screen.ListOrder = screen.ListOrder * -1;
		}

		screen.ListSort = col;
	}

	//
	// Toggle expanded.
	//
	viewToggle(screen: Screener) {
		if (screen.Expanded) {
			screen.Expanded = false;
		} else {
			screen.Expanded = true;
		}

		this.storePerferedView();
	}

	//
	// View change
	//
	viewChange(screen: Screener, type: string) {
		screen.View = type;
		this.storePerferedView();
	}

	//
	// Save our preferred views into local storage.
	//
	storePerferedView() {
		let obj = {}

		for (let i = 0; i < this.screeners.length; i++) {
			obj[this.screeners[i].Id] = { View: this.screeners[i].View, Expanded: this.screeners[i].Expanded };
		}

		localStorage.setItem("screener-view", JSON.stringify(obj));
	}

	//
	// Get and set preferred views.
	//
	getSetpreferedViews() {
		let json = localStorage.getItem("screener-view");

		if (!json) {
			for (let i = 0; i < this.screeners.length; i++) {
				this.screeners[i].View = "grid";
			}
			return;
		}

		let obj = JSON.parse(json);

		for (let i = 0; i < this.screeners.length; i++) {
			if (typeof obj[this.screeners[i].Id] != "undefined") {
				this.screeners[i].View = obj[this.screeners[i].Id].View;
				this.screeners[i].Expanded = obj[this.screeners[i].Id].Expanded;
			} else {
				this.screeners[i].View = "grid";
			}
		}

	}

	//
	// Get screeners.
	//
	getScreeners() {
		// Make api call to get screeners
		this.screenerService.get().subscribe(

			(res) => {

				// Assign the screeners.
				this.screeners = res;

				// Set the preferred views.
				this.getSetpreferedViews();

				// Make sure we have at least one screener
				if (this.screeners.length <= 0) {
					this.router.navigate(['/screener/add']);
					return;
				}

				// Load results.
				for (let i = 0; i < this.screeners.length; i++) {
					this.screenerService.getResults(this.screeners[i].Id, this.stateService.GetActiveBrokerAccount().Id).subscribe((res) => {
						this.screeners[i].Results = res;
					});
				}

				// Store in site cache.
				this.stateService.SetScreens(this.screeners);
			},

			// Error
			(err: HttpErrorResponse) => {

				if (err.error instanceof Error) {
					// A client-side or network error occurred. Handle it accordingly.
					console.log('An error occurred:', err.error);
				} else {
					if (err.error.error == 'No Record Found.') {
						this.router.navigate(['/screener/add'], { queryParams: { action: 'first-run' } });
					}
				}
			}
		);
	}

	//
	// Place trade from result
	//
	trade(screen: Screener, result: ScreenerResult) {

		// Just a double check
		if (result.Legs.length <= 0) {
			return;
		}

		// Set values
		let tradeDetails = new TradeDetails();
		tradeDetails.Symbol = result.Legs[0].OptionUnderlying;
		tradeDetails.Class = "multileg";

		if (result.Credit > 0) {
			tradeDetails.OrderType = "credit";
		} else {
			tradeDetails.OrderType = "debit";
		}

		// Set default duration.
		tradeDetails.Duration = "day";

		// Set default price
		tradeDetails.Price = result.MidPoint;

		// Build legs
		tradeDetails.Legs = [];

		// Get the default lot size
		let defaultLots = 1;

		switch (screen.Strategy) {
			case "put-credit-spread":
				defaultLots = this.settings.StrategyPcsLots;

				// If this is a mid point price.
				if (this.settings.StrategyPcsOpenPrice == "mid-point") {
					tradeDetails.Price = result.MidPoint;
				}

				// If this is a ask price.
				if (this.settings.StrategyPcsOpenPrice == "ask") {
					tradeDetails.Price = result.Credit;
				}

				// If this is a bid price.
				if (this.settings.StrategyPcsOpenPrice == "bid") {
					tradeDetails.Price = result.Debit;
				}
				break;
		}

		for (let i = 0; i < result.Legs.length; i++) {
			let side = "sell_to_close";
			let qty = defaultLots;

			// TODO this will need work based on the type of screener.
			if (i == 1) {
				side = "sell_to_open";
			} else {
				side = "buy_to_open"
			}

			tradeDetails.Legs.push(new TradeOptionLegs().createNew(result.Legs[i], result.Legs[i].OptionExpire, result.Legs[i].OptionType, result.Legs[i].OptionStrike, side, qty));
		}

		// Open builder to place trade.
		this.tradeService.tradeEvent.emit(new TradeEvent().createNew("toggle-trade-builder", tradeDetails));


	}

}

/* End File */
