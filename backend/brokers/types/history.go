package types

type History struct {
	Id          string
	BrokerId    string
	Type        string
	Date        string
	Amount      float64
	Symbol      string
	Commission  float64
	Description string
	Price       float64
	Quantity    int64
	TradeType   string
}

/*

{
	"history": {
		"event": [{
			"amount": -667.56,
			"date": "2017-04-03T00:00:00Z",
			"type": "trade",
			"trade": {
				"commission": 7.0000000000,
				"description": "PUT  SPY    05\/19\/17   223",
				"price": 1.100000,
				"quantity": 6.00000000,
				"symbol": "SPY170519P00223000",
				"trade_type": "Option"
			}
		}, {
			"amount": -110.09,
			"date": "2017-04-03T00:00:00Z",
			"type": "trade",
			"trade": {
				"commission": 0.0000000000,
				"description": "PUT  SPY    05\/19\/17   223",
				"price": 1.100000,
				"quantity": 1.00000000,
				"symbol": "SPY170519P00223000",
				"trade_type": "Option"
			}
		}, {
			"amount": 923.32,
			"date": "2017-04-03T00:00:00Z",
			"type": "trade",
			"trade": {
				"commission": 0.0000000000,
				"description": "PUT  SPY    05\/19\/17   225",
				"price": 1.320000,
				"quantity": -7.00000000,
				"symbol": "SPY170519P00225000",
				"trade_type": "Option"
			}
		}, {
			"amount": 41.34,
			"date": "2017-03-29T00:00:00Z",
			"type": "trade",
			"trade": {
				"commission": 7.0000000000,
				"description": "PUT  SPY    04\/07\/17   225",
				"price": 0.070000,
				"quantity": -7.00000000,
				"symbol": "SPY170407P00225000",
				"trade_type": "Option"
			}
		}, {
			"amount": -70.65,
			"date": "2017-03-29T00:00:00Z",
			"type": "trade",
			"trade": {
				"commission": 0.0000000000,
				"description": "PUT  SPY    04\/07\/17   227",
				"price": 0.100000,
				"quantity": 7.00000000,
				"symbol": "SPY170407P00227000",
				"trade_type": "Option"
			}
		}]
	}
}




vs.




{
	"history": {
		"event": {
			"amount": -667.56,
			"date": "2017-04-03T00:00:00Z",
			"type": "trade",
			"trade": {
				"commission": 7.0000000000,
				"description": "PUT  SPY    05\/19\/17   223",
				"price": 1.100000,
				"quantity": 6.00000000,
				"symbol": "SPY170519P00223000",
				"trade_type": "Option"
			}
		}
	}
}

*/
