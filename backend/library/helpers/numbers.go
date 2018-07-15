//
// Date: 2018-07-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-07-15
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

//
// Round
//
func Round(v float64, decimals int) float64 {
	var pow float64 = 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int((v*pow)+0.5)) / pow
}
