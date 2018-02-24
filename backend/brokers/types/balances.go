//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package types

type Balance struct {
	AccountNumber     string  `json:"account_number"`
	AccountValue      float64 `json:"account_value"`
	TotalCash         float64 `json:"total_cash"`
	OptionBuyingPower float64 `json:"option_buying_power"`
	StockBuyingPower  float64 `json:"stock_buying_power"`
}

/*

<accounts>
  <account>
    <balances>
      <account_number>12345678</account_number>
      <account_type>margin</account_type>
      <close_pl>0.000000</close_pl>
      <current_requirement>0.000000</current_requirement>
      <equity>2014.590000</equity>
      <long_market_value>0.000000</long_market_value>
      <market_value>0.000000</market_value>
      <open_pl>0.000000</open_pl>
      <option_long_value>0.000000</option_long_value>
      <option_requirement>0.000000</option_requirement>
      <option_short_value>0</option_short_value>
      <pending_orders_count>1</pending_orders_count>
      <short_market_value>0.000000</short_market_value>
      <stock_long_value>0.000000</stock_long_value>
      <total_cash>2014.590000</total_cash>
      <total_equity>2014.59000</total_equity>
      <uncleared_funds>0</uncleared_funds>
      <margin>
        <fed_call>0</fed_call>
        <maintenance_call>0</maintenance_call>
        <option_buying_power>1976.1000000000</option_buying_power>
        <stock_buying_power>3952.20000</stock_buying_power>
        <stock_short_value>0</stock_short_value>
        <sweep>0</sweep>
      </margin>
    </balances>
    <account_number>12345678</account_number>
  </account>
  <account>
    <balances>
      <account_number>87654321</account_number>
      <account_type>pdt</account_type>
      <close_pl>450.000000</close_pl>
      <current_requirement>133126.070000</current_requirement>
      <equity>228291.410000</equity>
      <long_market_value>214980.000000</long_market_value>
      <market_value>190460.000000</market_value>
      <open_pl>5651.920000</open_pl>
      <option_long_value>14400.000000</option_long_value>
      <option_requirement>32836.070000</option_requirement>
      <option_short_value>24520.0000000</option_short_value>
      <pending_orders_count>1</pending_orders_count>
      <short_market_value>24520.000000</short_market_value>
      <stock_long_value>200580.000000</stock_long_value>
      <total_cash>27711.410000</total_cash>
      <total_equity>218171.4100000000</total_equity>
      <uncleared_funds>0</uncleared_funds>
      <pdt>
        <day_trade_buying_power>569840.5100000000</day_trade_buying_power>
        <fed_call>0</fed_call>
        <maintenance_call>0</maintenance_call>
        <option_buying_power>95158.350000000001400000000</option_buying_power>
        <stock_buying_power>190316.7000000000028000</stock_buying_power>
        <stock_short_value>0</stock_short_value>
      </pdt>
    </balances>
    <account_number>87654321</account_number>
  </account>
</accounts>

*/
