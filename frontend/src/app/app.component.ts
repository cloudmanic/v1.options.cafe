import { Component } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';
import { Title } from '@angular/platform-browser';

declare let _paq: any;

@Component({
	selector: 'app-root',
	templateUrl: './app.component.html'
})

export class AppComponent {
	title = 'Options Cafe';

	//
	// Constructor
	//
	constructor(private router: Router, private titleService: Title) {
		// redirect
		let redt = localStorage.getItem('redirect');

		if (redt) {
			localStorage.removeItem('redirect');
			this.router.navigate([redt]);
		}

		// subscribe to router events and send page views to Analytics
		this.router.events.subscribe(event => {
			if (event instanceof NavigationEnd) {
				// We give it a timeout so we give time for the title to update.
				setTimeout(() => {
					_paq.push(['setCustomUrl', event.urlAfterRedirects]);
					_paq.push(['setDocumentTitle', this.titleService.getTitle()]);
					_paq.push(['setGenerationTimeMs', 0]);
					_paq.push(['trackPageView']);
					_paq.push(['enableLinkTracking']); // Should be at end.
				}, 50);
			}
		});
	}
}

/* End File */
