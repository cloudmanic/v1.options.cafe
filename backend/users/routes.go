//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import "github.com/tidwall/gjson"

// TODO: turn this in an a routes MAP type of code instead of a switch.
// Something like this....
// route.Add("refresh-watchlists", t.WsSendWatchlists)

//
// Listen for incoming feed requests.
//
func (t *Base) DoFeedRequestListen() {

	for {

		send := <-t.FeedRequestChan

		// Get message type
		msgType := gjson.Get(send.Message, "type").String()

		// Switch based on message type.
		switch msgType {

		// Refresh just the watchlists
		case "refresh-watchlists":
			t.WsSendWatchlists(t.Users[send.UserId], send.Message)
			break

		// Refresh all data from cache
		case "refresh-all-data":
			t.RefreshAllData(t.Users[send.UserId], send.Message)
			break

		}

	}

}

/* End File */
