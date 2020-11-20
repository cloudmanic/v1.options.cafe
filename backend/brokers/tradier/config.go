//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import "app.options.cafe/brokers/types"

//
// Return default config struct. Each brokers might have config values
// we track them here. One example is commissions.
//
func (t *Api) GetBrokerConfig() *types.BrokerConfig {
	return &types.BrokerConfig{
		DefaultStockCommission:   3.95,
		DefaultStockMin:          0.00,
		DefaultOptionCommission:  0.35,
		DefaultOptionSingleMin:   5.00,
		DefaultOptionMultiLegMin: 7.00,
		DefaultOptionBase:        0.00,
	}
}

/* End File */
