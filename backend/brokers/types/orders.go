package types

type OrderPreview struct {
	Status            string  `json:"status"`
	Error             string  `json:"error"`
	Commission        float64 `json:"commission"`
	Cost              float64 `json:"cost"`
	Fees              float64 `json:"fees"`
	Symbol            string  `json:"symbol"`
	Type              string  `json:"type"`
	Duration          string  `json:"duration"`
	Price             float64 `json:"price"`
	OrderCost         float64 `json:"order_cost"`
	MarginChange      float64 `json:"margin_change"`
	OptionRequirement float64 `json:"option_requirement"`
	Class             string  `json:"class"`
	Strategy          string  `json:"strategy"`
}

type Order struct {
	Id                string     `json:"-"`
	AccountId         string     `json:"account_id"`
	BrokerAccountId   uint       `json:"broker_account_id"`
	Type              string     `json:"type"`
	Symbol            string     `json:"symbol"`
	Side              string     `json:"side"`
	Quantity          float64    `json:"quantity"`
	Status            string     `json:"-"`
	Duration          string     `json:"duration"`
	Price             float64    `json:"price"`
	AvgFillPrice      float64    `json:"-"`
	ExecQuantity      float64    `json:"-"`
	LastFillPrice     float64    `json:"-"`
	LastFillQuantity  float64    `json:"-"`
	RemainingQuantity float64    `json:"-"`
	CreateDate        string     `json:"-"`
	TransactionDate   string     `json:"-"`
	Class             string     `json:"class"`
	OptionSymbol      string     `json:"option_symbol"`
	NumLegs           int        `json:"-"`
	Legs              []OrderLeg `json:"legs"`
}

type OrderLeg struct {
	Type              string  `json:"type"`
	Symbol            string  `json:"symbol"`
	OptionSymbol      string  `json:"option_symbol"`
	Side              string  `json:"side"`
	Quantity          float64 `json:"quantity"`
	Status            string  `json:"-"`
	Duration          string  `json:"duration"`
	AvgFillPrice      float64 `json:"-"`
	ExecQuantity      float64 `json:"-"`
	LastFillPrice     float64 `json:"-"`
	LastFillQuantity  float64 `json:"-"`
	RemainingQuantity float64 `json:"-"`
	CreateDate        string  `json:"-"`
	TransactionDate   string  `json:"-"`
}

// Tradier Temp Object for json decodeing
type TradierOrder struct {
	Id                int
	AccountId         string
	Type              string
	Symbol            string
	Side              string
	Quantity          float64
	Status            string
	Duration          string
	Price             float64
	AvgFillPrice      float64 `json:"avg_fill_price"`
	ExecQuantity      float64 `json:"exec_quantity"`
	LastFillPrice     float64 `json:"last_fill_price"`
	LastFillQuantity  float64 `json:"last_fill_quantity"`
	RemainingQuantity float64 `json:"remaining_quantity"`
	CreateDate        string  `json:"create_date"`
	TransactionDate   string  `json:"transaction_date"`
	Class             string
	OptionSymbol      string            `json:"option_symbol"`
	NumLegs           int               `json:"num_legs"`
	Legs              []TradierOrderLeg `json:"leg"`
}

type TradierOrderLeg struct {
	Type              string
	Symbol            string
	OptionSymbol      string `json:"option_symbol"`
	Side              string
	Quantity          float64
	Status            string
	Duration          string
	AvgFillPrice      float64 `json:"avg_fill_price"`
	ExecQuantity      float64 `json:"exec_quantity"`
	LastFillPrice     float64 `json:"last_fill_price"`
	LastFillQuantity  float64 `json:"last_fill_quantity"`
	RemainingQuantity float64 `json:"remaining_quantity"`
	CreateDate        string  `json:"create_date"`
	TransactionDate   string  `json:"transaction_date"`
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
