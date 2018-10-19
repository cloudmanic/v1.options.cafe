//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Settings } from '../../models/settings';
import { Component, OnInit } from '@angular/core';
import { SettingsService } from '../../providers/http/settings.service';
import { NotificationsService } from '../../providers/http/notifications.service';

declare var OneSignal: any;

@Component({
  selector: 'app-trading',
  templateUrl: './trading.component.html',
  styleUrls: []
})

export class TradingComponent implements OnInit 
{
  settings: Settings = new Settings();
  strategySettingsState: StrategyActiveState = new StrategyActiveState();

  //
  // Construct.
  //
  constructor(private notificationsService: NotificationsService, private settingsService: SettingsService) { }

  //
  // NgInit
  //
  ngOnInit() 
  {
    // Load data for page.
    this.loadSettingsData(); 

    //this.storeOneSignalUserId();


    // OneSignal.push(function() {
    //   OneSignal.isPushNotificationsEnabled().then(function(isEnabled) {
    //     if (isEnabled)
    //       console.log("Push notifications are enabled!");
    //     else
    //       console.log("Push notifications are not enabled yet.");
    //   });

    // });

  }

  //
  // Load settings data.
  //
  loadSettingsData()
  {
    this.settingsService.get().subscribe((res) => {
      this.settings = res;
      console.log(res);
    });
  }

  //
  // Change which setting we are on.
  // 
  strategySettingsClick(type: string)
  {
    // Clear all states.
    for(var row in this.strategySettingsState)
    {
      this.strategySettingsState[row] = false;
    }

    // Set the active state.
    this.strategySettingsState[type] = true;
  }

  //
  // Helper toggle for strategy 
  //
  strategyHelperToggle(type: string)
  {
    if(this.strategySettingsState[type])
    {
      this.strategySettingsState[type] = false;
    } else
    {
      this.strategySettingsState[type] = true;
    }
  }

  //
  // We call this when a Push Notification 
  // Checkbox is checked this will get 
  // approval from the user for push notifications
  // and subscribe the user to one signal.
  //
  setupBrowserNotifications()
  {

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
  storeOneSignalUserId()
  {
    // Get the user id and send it to the server.
    OneSignal.getUserId((userId) => {

      // Install channel
      this.notificationsService.createNotificationChannel('web-push', userId).subscribe();
 
    });

  }

}

//
// Keep track of what setting state we are in.
//
class StrategyActiveState
{
  PutCreditSpread: boolean = true;
  CallCreditSpread: boolean = false;
  PutDebitSpread: boolean = false;
  CallDebitSpread: boolean = false;

  // Helper Bubbles
  HelperPutCreditSpreadLots: boolean = false;
  HelperPutCreditSpreadOpenPrice: boolean = false;
  HelperPutCreditSpreadClosePrice: boolean = false;    
}
