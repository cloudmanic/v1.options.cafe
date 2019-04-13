//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { AnalyzeService } from '../../providers/http/analyze.service'
import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { WebsocketService } from '../../providers/http/websocket.service';
import { NotificationsService } from '../../providers/http/notifications.service';

@Component({
	selector: 'app-dashboard',
	templateUrl: './home.component.html'
})

export class DashboardComponent implements OnInit {
	openPrice: number = 88.99;
	showNotice: boolean = false;
	noticeId: number = 0;
	noticeTitle: string = '';
	noticeBody: string = '';
	ws_reconnecting: boolean = false;

	//
	// Construct...
	//
	constructor(private notificationsService: NotificationsService, private websocketService: WebsocketService, private changeDetect: ChangeDetectorRef) {
		// Load data for page.
		this.getNotifications();
	}

	//
	// On Init...
	//
	ngOnInit() {
		// Subscribe to when we are reconnecting to a websocket - Core
		this.websocketService.wsReconnecting.subscribe(data => {
			this.ws_reconnecting = data;
			this.changeDetect.detectChanges();
		});

	}

	//
	// See if we have any notifications to display.
	//
	getNotifications() {
		this.notificationsService.get("in-app", "dashboard-notice", "pending").subscribe(data => {

			// We only take one notice at a time.
			if (data.length <= 0) {
				this.noticeId = 0;
				this.showNotice = false;
				this.noticeTitle = "";
				this.noticeBody = "";
				return;
			}

			// Grab the first notice and display.
			this.noticeId = data[0].Id;
			this.noticeTitle = data[0].Title;
			this.noticeBody = data[0].LongMessage;
			this.showNotice = true;

		});
	}

	//
	// Close Notice
	//
	closeNotice() {
		// Send API call to mark notice as seen.
		this.notificationsService.markSeen(this.noticeId).subscribe(data => {

			// See if there is a "next" notice.
			this.getNotifications();

		});
	}

}

/* End File */
