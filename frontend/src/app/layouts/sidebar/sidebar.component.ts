//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Rx';
import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { Component, OnInit } from '@angular/core';
import { Balance } from '../../models/balance';
import { ChangeDetected } from '../../models/change-detected';
import { MarketStatus } from '../../models/market-status';
import { Broker } from '../../models/broker';
import { BrokerAccount } from '../../models/broker-account';
import { BrokerService } from '../../providers/http/broker.service';
import { StateService } from '../../providers/state/state.service';
import { StatusService } from '../../providers/http/status.service';
import { WebsocketService } from '../../providers/http/websocket.service';
import { MeService } from 'app/providers/http/me.service';

@Component({
	selector: 'app-sidebar',
	templateUrl: './sidebar.component.html'
})
export class SidebarComponent implements OnInit {

	brokerList: Broker[];
	balance: Balance;
	marketStatus: MarketStatus;
	selectedAccount: BrokerAccount;
	brokerAccountList: BrokerAccount[];

	private destory: Subject<boolean> = new Subject<boolean>();

	//
	// Construct.
	//
	constructor(private websocketService: WebsocketService, private brokerService: BrokerService, private stateService: StateService, private statusService: StatusService, private meService: MeService) {
		// Load data to start
		this.brokerAccountList = [];
		this.getBrokers();
	}

	//
	// Oninit...
	//
	ngOnInit() {
		// Subscribe to changes in the selected broker.
		this.stateService.BrokerChange.takeUntil(this.destory).subscribe(data => {
			this.getBalances();
		});

		// Subscribe to data updates from the broker - Balances
		this.websocketService.balancesPush.subscribe(data => {
			this.doBalanaces(data);
		});

		// Subscribe to data updates from the broker - Market Status
		this.websocketService.marketStatusPush.subscribe(data => {
			this.doMarketStatus(data);
		});

		// Subscribe to when changes are detected at the server.
		this.websocketService.changedDetectedPush.takeUntil(this.destory).subscribe(data => {
			this.manageChangeDetection(data);
		});

		// Get the user profile so we can store it in local storage.
		this.meService.getProfile().subscribe((res) => {
			localStorage.setItem('user_email', res.Email);
		});
	}

	//
	// OnDestroy
	//
	ngOnDestroy() {
		this.destory.next();
		this.destory.complete();
	}

	//
	// Manage change detection.
	//
	manageChangeDetection(data: ChangeDetected) {
		if (data.Type == 'brokers') {
			this.getBrokers();
		}
	}

	//
	// Do market status.
	//
	doMarketStatus(status: MarketStatus) {

		// Only update on change
		if (JSON.stringify(status) != JSON.stringify(this.marketStatus)) {
			this.marketStatus = status;
		}

	}

	//
	// Get balances
	//
	getBalances() {
		if (!this.selectedAccount) {
			return
		}

		// Get balance data
		this.brokerService.getAccountBalance(this.selectedAccount.BrokerId, this.selectedAccount.Id).subscribe((data) => {
			this.balance = data;
		});
	}

	//
	// Get brokers
	//
	getBrokers() {
		// AJAX call to get brokers
		this.brokerService.get().subscribe((data) => {

			this.brokerList = data;

			let activeAccountId: string = this.stateService.GetStoredActiveAccountId();

			// Default to first one.
			if (!activeAccountId) {
				// Make sure we have at least one broker.
				if (!this.brokerList[0]) {
					return false;
				}

				// Do we have a stored broker
				this.stateService.SetActiveBrokerAccount(this.brokerList[0].BrokerAccounts[0]);
				activeAccountId = this.stateService.GetStoredActiveAccountId();
			}

			this.brokerAccountList = [];

			// Loop through all the brokers and set our active broker. And make a list.
			for (var k = 0; k < this.brokerList.length; k++) {
				for (var i = 0; i < this.brokerList[k].BrokerAccounts.length; i++) {
					this.brokerAccountList.push(this.brokerList[k].BrokerAccounts[i]);

					// Set the selected account.
					if (String(this.brokerList[k].BrokerAccounts[i].Id) == activeAccountId) {
						this.selectedAccount = this.brokerList[k].BrokerAccounts[i];

						// Force refresh of balances
						this.stateService.SetActiveBrokerAccount(this.selectedAccount);
					}
				}
			}
		});
	}

	//
	// Do balances
	//
	doBalanaces(data: Balance[]) {
		// We have not gotten our brokers yet.
		if (!this.stateService.GetActiveBrokerAccount()) {
			return false;
		}

		for (var i = 0; i < data.length; i++) {
			if (data[i].AccountNumber == this.stateService.GetActiveBrokerAccount().AccountNumber) {
				this.balance = data[i];
			}
		}
	}

	//
	// On account change.
	//
	onAccountChange() {
		this.stateService.SetActiveBrokerAccount(this.selectedAccount);
	}
}

/* End File */
