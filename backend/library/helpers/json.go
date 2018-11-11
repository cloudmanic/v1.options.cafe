//
// Date: 2018-11-11
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-11
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import "encoding/json"

//
// Json encode no error
//
func JsonEncode(obj interface{}) string {
	rawJson, err := json.Marshal(obj)

	if err != nil {
		return "{}"
	}

	// Return happy
	return string(rawJson)
}

/* End File */
