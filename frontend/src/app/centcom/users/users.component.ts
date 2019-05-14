//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Observable';
import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Component, OnInit } from '@angular/core';

@Component({
	selector: 'app-users',
	templateUrl: './users.component.html',
	styleUrls: []
})

export class UsersComponent implements OnInit {
	users: Object[] = [];

	//
	// Constructor
	//
	constructor(private http: HttpClient) { }

	//
	// NgInit
	//
	ngOnInit() {
		this.getUsers();
	}

	//
	// Get users
	//
	getUsers() {
		this.http.get<User[]>(environment.app_server + '/api/admin/users').subscribe((data) => {
			this.users = data;
			console.log(data);
		});
	}

	//
	// Delete a user.
	//
	deleteUser(user: User) {
		let c = confirm("Are you sure you want to delete " + user.email + "? MAKE SURE YOU TOOK A DB BACKUP!!!!");

		if (!c) {
			return
		}

		// Send request to delete the user for good.
		this.http.delete<boolean>(environment.app_server + '/api/admin/users/' + user.id).subscribe(

			(data) => {
				this.getUsers();
				return true;
			},

			// Error
			(err: HttpErrorResponse) => {
				alert(err);
			}

		);
	}

	//
	// Login as user.
	//
	loginAsUser(user: User) {
		let post = {
			id: user.id
		}

		this.http.post<Object>(environment.app_server + '/api/admin/users/login-as-user', post).subscribe((data) => {

			// Store old values
			localStorage.setItem('user_id_centcom', localStorage.getItem('user_id'));
			localStorage.setItem('access_token_centcom', localStorage.getItem('access_token'));
			localStorage.setItem('active_account_centcom', localStorage.getItem('active_account'));
			localStorage.setItem('active_watchlist_centcom', localStorage.getItem('active_watchlist'));

			// Remove local storage
			localStorage.removeItem('user_id');
			localStorage.removeItem('redirect');
			localStorage.removeItem('broker_new_id');
			localStorage.removeItem('access_token');
			localStorage.removeItem('active_account');
			localStorage.removeItem('active_watchlist');

			// Store access token in local storage.
			localStorage.setItem('user_id', data["user_id"].toString());
			localStorage.setItem('access_token', data["access_token"]);

			// Redirect to app
			window.location.href = '/';
		});
	}

}

interface User {
	id: number;
	first_name: string;
	last_name: string;
	email: string;
	last_activity: string;
}

/* End File */
