//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Me } from '../../models/me';
import { Settings } from '../../models/settings';
import { HttpErrors } from '../../models/http-errors';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { MeService } from '../../providers/http/me.service';
import { StateService } from '../../providers/state/state.service';
import { SettingsService } from '../../providers/http/settings.service';
import { NotificationsService } from '../../providers/http/notifications.service';
import { environment } from 'environments/environment';
import { Title } from '@angular/platform-browser';

declare var OneSignal: any;

const pageTitle: string = environment.title_prefix + "Settings Trading";

@Component({
	selector: 'app-trading',
	templateUrl: './trading.component.html',
	styleUrls: []
})

export class TradingComponent implements OnInit {
	tmpUserProfile: Me = new Me();
	userProfile: Me = new Me();
	showEditPhone: boolean = false;
	settings: Settings = new Settings();
	httpErrors: HttpErrors = new HttpErrors();
	strategySettingsState: StrategyActiveState = new StrategyActiveState();

	//
	// Construct.
	//
	constructor(private notificationsService: NotificationsService, private settingsService: SettingsService, private stateService: StateService, private meService: MeService, private titleService: Title) {
		this.settings = this.stateService.GetSettings();
		this.userProfile = this.stateService.GetSettingsUserProfile();
	}

	//
	// NgInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Load data for page.
		this.loadUserProfile();
		this.loadSettingsData();
	}

	//
	// Get user data.
	//
	loadUserProfile() {
		this.meService.getProfile().subscribe((res) => {
			this.userProfile = res;
			this.stateService.SetSettingsUserProfile(res);
		});
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
	// Update settings
	//
	updateSettings() {
		this.settingsService.update(this.settings).subscribe((res) => {
			this.stateService.SetSettings(this.settings);

			// Show success notice
			this.stateService.SiteSuccess.emit("Your trade settings have been successfully updated.");
		});
	}

	//
	// Strategy change
	//
	strategyChange() {
		this.updateSettings();
	}

	//
	// Notice change
	//
	noticeChange(which: string) {
		if (this.settings[which] == "Yes") {
			this.settings[which] = "No";
		} else {
			this.settings[which] = "Yes";
		}

		// is this a SMS call?
		if ((which.indexOf("Sms") > 0) && (this.userProfile.Phone.length <= 0)) {
			this.tmpUserProfile = new Me().setFromObj(this.userProfile);
			this.showEditPhone = true;
			return;
		}

		// is this a Push call?
		if (which.indexOf("Push") > 0) {
			this.setupBrowserNotifications();
		}


		this.updateSettings();
	}

	//
	// Change which setting we are on.
	//
	strategySettingsClick(type: string) {
		// Clear all states.
		for (var row in this.strategySettingsState) {
			this.strategySettingsState[row] = false;
		}

		// Set the active state.
		this.strategySettingsState[type] = true;
	}

	//
	// Helper toggle for strategy
	//
	strategyHelperToggle(type: string) {
		if (this.strategySettingsState[type]) {
			this.strategySettingsState[type] = false;
		} else {
			this.strategySettingsState[type] = true;
		}
	}

	//
	// We call this when a Push Notification
	// Checkbox is checked this will get
	// approval from the user for push notifications
	// and subscribe the user to one signal.
	//
	setupBrowserNotifications() {

		OneSignal.push(() => {

			OneSignal.registerForPushNotifications();
			OneSignal.setSubscription(true);

			// Tag this user at One Signal
			let userId = localStorage.getItem('user_id');
			if (userId.length) {
				OneSignal.sendTags({ userId: userId });
			}

			// Send to server
			this.storeOneSignalUserId();
		});

	}

	//
	// Send OneSignal User ID to backend
	//
	storeOneSignalUserId() {
		// Get the user id and send it to the server.
		OneSignal.getUserId((userId) => {

			// Install channel
			this.notificationsService.createNotificationChannel('web-push', userId).subscribe();

		});

	}

	//
	// Save Edit profile.
	//
	doSaveEditProfile() {
		let yesError = false;

		// Clear Validation
		this.httpErrors = new HttpErrors();

		// Validate - First Name
		if (this.userProfile.FirstName.length <= 0) {
			yesError = true;
			this.httpErrors.FirstName = "First name field is required.";
		}

		// Validate - Last Name
		if (this.userProfile.LastName.length <= 0) {
			yesError = true;
			this.httpErrors.LastName = "Last name field is required.";
		}

		// Validate - Email
		if (this.userProfile.Email.length <= 0) {
			yesError = true;
			this.httpErrors.Email = "Email field is required.";
		}

		if (yesError) {
			return;
		}

		// Ajax call to save the profile data.
		this.meService.saveProfile(this.userProfile).subscribe(

			// Success
			(res) => {
				this.userProfile = res;
				this.showEditPhone = false;
				this.stateService.SetSettingsUserProfile(res);

				this.updateSettings();

				// Show success notice
				this.stateService.SiteSuccess.emit("Your profile has been successfully updated.");
			},

			// Error
			(err: HttpErrorResponse) => {
				this.httpErrors = new HttpErrors().fromJson(err.error.errors);
			}

		);
	}

	//
	// Cancel profile.
	//
	doCancelEditProfile() {
		// Clear Validation
		this.httpErrors = new HttpErrors();

		// Reset values
		this.userProfile = new Me().setFromObj(this.tmpUserProfile);
		this.showEditPhone = false;

		// Reset values
		this.settings.NoticeMarketOpenedSms = "No";
		this.settings.NoticeMarketClosedSms = "No";
		this.settings.NoticeTradeFilledSms = "No";
	}

	//
	// Do close tool top.
	//
	closeToolTip(agreed: boolean) {
		this.strategySettingsState.HelperVerticalSpreadLots = false;
		this.strategySettingsState.HelperVerticalSpreadOpenPrice = false;
		this.strategySettingsState.HelperVerticalSpreadClosePrice = false;
	}

}

//
// Keep track of what setting state we are in.
//
class StrategyActiveState {
	VerticalSpread: boolean = true;

	// Helper Bubbles
	HelperVerticalSpreadLots: boolean = false;
	HelperVerticalSpreadOpenPrice: boolean = false;
	HelperVerticalSpreadClosePrice: boolean = false;


	// Tool tips
	HelperVerticalSpreadOpenTitle: string = "Open Price";
	HelperVerticalSpreadOpenMsg: string = "When opening a trade from the Options Cafe screener you can choose if you want to open at the Bid, Ask, or Midpoint price. This is simply the default setting -- you can override your open price at the time of your trade.";

	HelperVerticalSpreadLotsTitle: string = "Trade Lots";
	HelperVerticalSpreadLotsMsg: string = "When opening a trade from the Options Cafe screener you may have a standard number of contracts (or lots) you want to trade. You can set your default lots size here for quick ordering. Of course at the time you place your trade you can always change how many contracts you order.";

	HelperVerticalSpreadCloseTitle: string = "Close Price";
	HelperVerticalSpreadCloseMsg: string = "When closing a position many options traders have a desired price. For example you might open a put credit spread at a credit of $0.20 and want to close it at a debit of $0.03. With this setting you can set your default close price. This feature is just for one click trading. You have to place the close order yourself. Options Cafe will not auto close trades for you.";
}

/* End File */
