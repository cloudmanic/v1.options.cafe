package types

type Order struct {
  Id int
  AccountId string  
  Type string
  Symbol string
  Side string
  Quantity float64
  Status string
  Duration string
  Price float64
  AvgFillPrice float64 
  ExecQuantity float64 
  LastFillPrice float64
  LastFillQuantity float64
  RemainingQuantity float64
  CreateDate string
  TransactionDate string
  Class string
  NumLegs int
  Legs []OrderLeg   
}

type OrderLeg struct {
  Type string
  Symbol string
  OptionSymbol string
  Side string
  Quantity float64
  Status string
  Duration string
  AvgFillPrice float64
  ExecQuantity float64
  LastFillPrice float64
  LastFillQuantity float64
  RemainingQuantity float64
  CreateDate string
  TransactionDate string  
}

// Tradier Temp Object for json decodeing
type TradierOrder struct {
  Id int
  AccountId string   
  Type string
  Symbol string
  Side string
  Quantity float64
  Status string
  Duration string
  Price float64
  AvgFillPrice float64 `json:"avg_fill_price"`
  ExecQuantity float64 `json:"exec_quantity"`
  LastFillPrice float64 `json:"last_fill_price"`
  LastFillQuantity float64 `json:"last_fill_quantity"`
  RemainingQuantity float64 `json:"remaining_quantity"`
  CreateDate string `json:"create_date"`
  TransactionDate string `json:"transaction_date"`
  Class string
  NumLegs int `json:"num_legs"`
  Legs []TradierOrderLeg `json:"leg"`   
}

type TradierOrderLeg struct {
  Type string
  Symbol string
  OptionSymbol string `json:"option_symbol"`
  Side string
  Quantity float64
  Status string
  Duration string
  AvgFillPrice float64 `json:"avg_fill_price"`
  ExecQuantity float64 `json:"exec_quantity"`
  LastFillPrice float64 `json:"last_fill_price"`
  LastFillQuantity float64 `json:"last_fill_quantity"`
  RemainingQuantity float64 `json:"remaining_quantity"`
  CreateDate string `json:"create_date"`
  TransactionDate string `json:"transaction_date"`  
}

/*
<orders>
  <order>
    <id>5751</id>
    <type>limit</type>
    <symbol>F</symbol>
    <side>sell_short</side>
    <quantity>1.00000</quantity>
    <status>pending</status>
    <duration>day</duration>
    <price>21.16000</price>
    <avg_fill_price>0.00000</avg_fill_price>
    <exec_quantity>0.00000</exec_quantity>
    <last_fill_price>0.00000</last_fill_price>
    <last_fill_quantity>0.00000</last_fill_quantity>
    <remaining_quantity>0.00000</remaining_quantity>
    <create_date>2014-05-28T12:04:52.627Z</create_date>
    <transaction_date>2014-05-28T12:04:52.627Z</transaction_date>
    <class>equity</class>
  </order>
  <order>
    <id>5776</id>
    <type>market</type>
    <symbol>F</symbol>
    <side>buy</side>
    <quantity>0.00000</quantity>
    <status>pending</status>
    <duration>day</duration>
    <avg_fill_price>0.00000</avg_fill_price>
    <exec_quantity>0.00000</exec_quantity>
    <last_fill_price>0.00000</last_fill_price>
    <last_fill_quantity>0.00000</last_fill_quantity>
    <remaining_quantity>0.00000</remaining_quantity>
    <create_date>2014-05-28T12:05:51.673Z</create_date>
    <transaction_date>2014-05-28T12:05:51.673Z</transaction_date>
    <class>combo</class>
    <num_legs>2</num_legs>
    <leg>
      <type>market</type>
      <symbol>F</symbol>
      <side>buy</side>
      <quantity>100.00000</quantity>
      <status>pending</status>
      <duration>day</duration>
      <avg_fill_price>0.00000</avg_fill_price>
      <exec_quantity>0.00000</exec_quantity>
      <last_fill_price>0.00000</last_fill_price>
      <last_fill_quantity>0.00000</last_fill_quantity>
      <remaining_quantity>0.00000</remaining_quantity>
      <create_date>2014-05-28T12:05:51.660Z</create_date>
      <transaction_date>2014-05-28T12:05:51.660Z</transaction_date>
      <option_symbol></option_symbol>
    </leg>
    <leg>
      <type>market</type>
      <symbol>F</symbol>
      <side>sell_to_open</side>
      <quantity>1.00000</quantity>
      <status>pending</status>
      <duration>day</duration>
      <avg_fill_price>0.00000</avg_fill_price>
      <exec_quantity>0.00000</exec_quantity>
      <last_fill_price>0.00000</last_fill_price>
      <last_fill_quantity>0.00000</last_fill_quantity>
      <remaining_quantity>0.00000</remaining_quantity>
      <create_date>2014-05-28T12:05:51.673Z</create_date>
      <transaction_date>2014-05-28T12:05:51.673Z</transaction_date>
      <option_symbol>F140530C00016000</option_symbol>
    </leg>
  </order>
</orders>
*/