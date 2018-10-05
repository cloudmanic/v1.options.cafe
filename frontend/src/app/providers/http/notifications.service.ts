//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Observable';
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

declare var OneSignal: any;

@Injectable()
export class NotificationsService 
{
  //
  // Construct.
  //
  constructor(private http: HttpClient) 
  { 
    this.setupOpenSignal();
  }

  //
  // Setup open signal
  //
  setupOpenSignal()
  {
    // Setup OneSignal
    OneSignal.push(["init", {
      appId: environment.one_signal_app_id
    }]);
  }

  //
  // Get brokers
  //
  createNotificationChannel(type: string, channel_id: string): Observable<boolean> {
    return this.http.post<boolean>(environment.app_server + '/api/v1/notifications/add-channel', { type: type, channel_id: channel_id}).map(
      (data) => {
        return true
      });
  }   

}

/* End File */