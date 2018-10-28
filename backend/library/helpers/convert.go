//
// Date: 3/5/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"strconv"
)

//
// Float to String
func FloatToString(fv float64) string {
	return strconv.FormatFloat(fv, 'f', 2, 64)
}

//
// Convert a string to a int
//
func StringToInt(s string) int {
	idInt, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return idInt
}

//
// Convert a string to a uint
//
func StringToUint(s string) uint {
	idInt, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return uint(idInt)
}

/* End File */
