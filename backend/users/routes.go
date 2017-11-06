//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import (
	"app.options.cafe/backend/controllers"
	"github.com/tidwall/gjson"
)

//
// Get Routes
//
func (t *Base) GetRoutes() map[string]func(*UserFeed, controllers.ReceivedStruct) {

	routes := make(map[string]func(*UserFeed, controllers.ReceivedStruct))

	// Set routes
	routes["refresh-watchlists"] = t.WsSendWatchlists
	routes["refresh-all-data"] = t.RefreshAllData

	// Return happy
	return routes
}

//
// Listen for incoming feed requests.
//
func (t *Base) DoFeedRequestListen() {

	routes := t.GetRoutes()

	for {

		send := <-t.FeedRequestChan

		// Get message type
		msgType := gjson.Get(send.Message, "type").String()

		// Make sure we know about this type & Call function to manage request
		if _, ok := routes[msgType]; ok {
			routes[msgType](t.Users[send.UserId], send)
		}
	}

}

/* End File */
