//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { NotificationsService } from '../../providers/http/notifications.service';

declare var OneSignal: any;

@Component({
  selector: 'app-trading',
  templateUrl: './trading.component.html',
  styleUrls: []
})

export class TradingComponent implements OnInit {


  //
  // Construct.
  //
  constructor(private notificationsService: NotificationsService) { }

  //
  // NgInit
  //
  ngOnInit() 
  {
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
      this.notificationsService.createNotificationChannel('One Signal', userId).subscribe();
 
    });

  }

}
