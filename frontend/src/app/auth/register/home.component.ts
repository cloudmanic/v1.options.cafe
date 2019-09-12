//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { NgForm } from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router, ActivatedRoute, Params } from '@angular/router';
import { environment } from '../../../environments/environment';
import { Title } from '@angular/platform-browser';

interface RegisterResponse {
	status: number,
	user_id: number,
	error: string,
	access_token: string
}

interface LoginResponse {
	status: number,
	user_id: number,
	error: string,
	access_token: string,
	broker_count: number
}

interface GoogleSessionResponse {
	session_secret: string
}

const pageTitle: string = environment.title_prefix + "Register";

@Component({
	selector: 'app-auth-register',
	templateUrl: './home.component.html'
})

export class AuthRegisterComponent implements OnInit {
	errorMsg = "";
	submitBtn = "Create Account";
	googleLoginState: boolean = false;

	// Form
	email = "";

	//
	// Constructor.
	//
	constructor(private http: HttpClient, private router: Router, private activatedRoute: ActivatedRoute, private titleService: Title) { }

	//
	// NgInit.
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// subscribe to router event
		this.activatedRoute.queryParams.subscribe((params: Params) => {

			// See if we have an email address already
			if (params['email']) {
				this.email = params['email'];
			}

			// See if we have a google_auth_success for a successful google login.
			if (params['google_auth_success']) {
				this.processPostGoogleLogin();
			}

			// See if our google login failed to login.
			if (params['google_auth_failed']) {
				if (params['google_auth_failed'] == 'user-already-in-system') {
					this.errorMsg = "Looks like you already have an account. Please try to login.";
				} else {
					this.errorMsg = "Looks like there was an error logging into Google.";
				}
			}
		});
	}

	//
	// Process post Google login.
	//
	processPostGoogleLogin() {
		let sessionKey = localStorage.getItem('google_auth_session_key');
		let sessionSecret = localStorage.getItem('google_auth_session_secret');

		// Remove keys we do not need.
		localStorage.removeItem('google_auth_session_key');
		localStorage.removeItem('google_auth_session_secret');


		// Make the the HTTP request:
		this.http.post<LoginResponse>(environment.app_server + '/oauth/google/token', { session_key: sessionKey, session_secret: sessionSecret, grant_type: "password", client_id: environment.client_id }).subscribe(

			// Success
			data => {
				// Store access token in local storage.
				localStorage.setItem('user_id', data.user_id.toString());
				localStorage.setItem('access_token', data.access_token);

				// Redirect to broker select
				this.router.navigate(['/broker-select']);
			},

			// Error
			(err: HttpErrorResponse) => {
				this.errorMsg = "Unable to login via your Google account.";
				console.log('An error occurred:', err);
			}

		);
	}

	//
	// Do Google login.
	//
	googleLogin() {
		this.googleLoginState = true;
		let sessionKey = this.getRandomString();
		localStorage.setItem('google_auth_session_key', sessionKey);

		// Make the the HTTP request:
		this.http.post<GoogleSessionResponse>(environment.app_server + '/oauth/google/session', { session_key: sessionKey, type: 'register', redirect: environment.site_url + '/register' }).subscribe(

			// Success
			data => {
				// Store session secret.
				localStorage.setItem('google_auth_session_secret', data.session_secret.toString());

				// Redirect to start the Google login
				window.location.href = environment.app_server + '/oauth/google?session_key=' + sessionKey;
			},

			// Error
			(err: HttpErrorResponse) => {
				console.log('An error occurred:', err);
			}

		);
	}

	//
	// Register submit.
	//
	onSubmit(form: NgForm) {

		if (this.googleLoginState) {
			return;
		}

		// Clear post error.
		this.errorMsg = "";

		// Update submit button
		this.submitBtn = "Saving...";

		// Make the the HTTP request:
		this.http.post<RegisterResponse>(environment.app_server + '/register', form.value).subscribe(

			// Success
			data => {

				// Store access token in local storage.
				localStorage.setItem('user_id', data.user_id.toString());
				localStorage.setItem('access_token', data.access_token);

				// Redirect to broker select
				this.router.navigate(['/broker-select']);

			},

			// Error
			(err: HttpErrorResponse) => {

				// Change button back.
				this.submitBtn = "Create Account";

				if (err.error instanceof Error) {
					// A client-side or network error occurred. Handle it accordingly.
					console.log('An error occurred:', err.error);
				} else {
					// Print error message
					this.errorMsg = err.error.error;
				}

			}

		);

	}

	//
	// Make a random string to manage the google access token request.
	//
	getRandomString() {
		let text = "";
		let possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

		for (var i = 0; i < 20; i++) {
			text += possible.charAt(Math.floor(Math.random() * possible.length));
		}

		return text;
	}

}
