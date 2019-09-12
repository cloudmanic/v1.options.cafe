import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { environment } from 'environments/environment';

const pageTitle: string = environment.title_prefix + "Settings Account";

@Component({
	selector: 'app-account',
	templateUrl: './account.component.html',
	styleUrls: []
})
export class AccountComponent implements OnInit {

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
	}

}
