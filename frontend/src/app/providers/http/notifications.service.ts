//
// Date: 10/4/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Observable } from 'rxjs/Observable';
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Notification } from '../../models/notification';
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
  // Get Notifications
  //
  get(channel: string, uri: string, status: string): Observable<Notification[]> {
    let queryStr = '?order=id&sort=asc&status=' + status + '&channel=' + channel + '&uri=' + uri;

    return this.http.get<Notification[]>(environment.app_server + '/api/v1/notifications' + queryStr).map(
      (data) => { return new Notification().fromJsonList(data); });
  }

  //
  // Mark a notice as seen
  //
  markSeen(id: number): Observable<boolean> {
    return this.http.put<boolean>(environment.app_server + '/api/v1/notifications/' +  id, { status: "seen" }).map(
      (data) => {
        return true
      });
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