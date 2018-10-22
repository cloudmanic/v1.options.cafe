//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';

export class Notification 
{
  public Id: number;
  public Status: string;
  public Channel: string;
  public Uri: string;
  public UriRefId: string;
  public Title: string;
  public ShortMessage: string;
  public LongMessage: string;
  public SentTime: Date;

  //
  // Build from json list.
  //
  fromJsonList(json: Object[]): Notification[] {
    let list: Notification[] = [];

    if (!json) 
    {
      return list;
    }

    for (let i = 0; i < json.length; i++) 
    {
      list.push(this.fromJson(json[i]));
    }

    return list;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): Notification 
  {
    let obj = new Notification();

    obj.Id = json["id"];
    obj.Status = json["status"];
    obj.Channel = json["channel"];
    obj.Uri = json["uri"];
    obj.UriRefId = json["uri_ref_id"];
    obj.Title = json["title"];
    obj.ShortMessage = json["short_message"];
    obj.LongMessage = json["long_message"];
    obj.SentTime = moment(json["sent_time"]).toDate();

    return obj;
  }

  
}

/* End File */