import { Component, OnInit } from '@angular/core';
import { environment } from 'environments/environment';
import { Title } from '@angular/platform-browser';

declare var TradingView: any;

const pageTitle: string = environment.title_prefix + "Research Symbol";

@Component({
	selector: 'app-research-symbol',
	templateUrl: './symbol.component.html',
	styleUrls: []
})

export class SymbolComponent implements OnInit {
	//
	// Constructor
	//
	constructor(private titleService: Title) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Load Trading chart
		new TradingView.widget(
			{
				"width": "100%",
				"height": 700,
				"symbol": "AMEX:SPY",
				"interval": "D",
				"timezone": "America/Los_Angeles",
				"theme": "Light",
				"style": "1",
				"locale": "en",
				"toolbar_bg": "#f1f3f6",
				"enable_publishing": false,
				"allow_symbol_change": true,
				"container_id": "tradingview_300fb"
			}
		);


	}

}

/* End File */
