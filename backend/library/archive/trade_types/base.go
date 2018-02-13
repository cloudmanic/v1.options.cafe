//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package trade_types

//
// Find max strike price
//
func FindMaxStike(strikes []float64) float64 {

	max := 0.00

	// Find max strike
	for _, e := range strikes {
		if e > max {
			max = e
		}
	}

	return max
}

//
// Find min strike price
//
func FindMinStike(strikes []float64) float64 {

	min := strikes[0]

	// Find min strike
	for _, e := range strikes {
		if e < min {
			min = e
		}
	}

	return min
}

/* End File */
