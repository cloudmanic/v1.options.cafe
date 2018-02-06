package types

type Position struct {
	Id           int
	AccountId    string
	Symbol       string
	DateAcquired string `json:date_acquired`
	CostBasis    float64
	Quantity     float64
}

/*
{
	"accounts": {
		"account": [{
			"account_number": "6YA06984",
			"positions": "null"
		}, {
			"account_number": "6YA05782",
			"positions": {
				"position": [{
					"cost_basis": 519.5,
					"date_acquired": "2018-01-09T11:15:51.981Z",
					"id": 266353,
					"quantity": 10.0,
					"symbol": "CONE"
				}, {
					"cost_basis": 387.2796,
					"date_acquired": "2018-01-09T11:15:51.999Z",
					"id": 266354,
					"quantity": 4.0,
					"symbol": "DIS"
				}, {
					"cost_basis": 395.7996,
					"date_acquired": "2018-01-09T11:15:52.016Z",
					"id": 266355,
					"quantity": 4.0,
					"symbol": "NFLX"
				}, {
					"cost_basis": 169.598,
					"date_acquired": "2018-01-09T11:15:52.034Z",
					"id": 266356,
					"quantity": 20.0,
					"symbol": "USO"
				}, {
					"cost_basis": 484.8396,
					"date_acquired": "2018-01-09T11:15:52.052Z",
					"id": 266357,
					"quantity": 4.0,
					"symbol": "VTI"
				}]
			}
		}, {
			"account_number": "6YA06085",
			"positions": {
				"position": [{
					"cost_basis": 1404.0,
					"date_acquired": "2018-02-05T16:17:18.308Z",
					"id": 286040,
					"quantity": 9.0,
					"symbol": "SPY180309P00262000"
				}, {
					"cost_basis": -1629.0,
					"date_acquired": "2018-02-05T16:17:18.308Z",
					"id": 286039,
					"quantity": -9.0,
					"symbol": "SPY180309P00264000"
				}, {
					"cost_basis": 2303.01,
					"date_acquired": "2018-02-05T20:14:14.167Z",
					"id": 286459,
					"quantity": 9.0,
					"symbol": "SPY180316P00250000"
				}, {
					"cost_basis": 2097.0,
					"date_acquired": "2018-02-05T19:26:06.234Z",
					"id": 286357,
					"quantity": 9.0,
					"symbol": "SPY180316P00251000"
				}, {
					"cost_basis": -2618.01,
					"date_acquired": "2018-02-05T20:14:14.167Z",
					"id": 286458,
					"quantity": -9.0,
					"symbol": "SPY180316P00252000"
				}, {
					"cost_basis": -2295.0,
					"date_acquired": "2018-02-05T19:26:06.234Z",
					"id": 286356,
					"quantity": -9.0,
					"symbol": "SPY180316P00253000"
				}, {
					"cost_basis": 1332.0,
					"date_acquired": "2018-02-02T16:19:15.693Z",
					"id": 284951,
					"quantity": 9.0,
					"symbol": "SPY180316P00264000"
				}, {
					"cost_basis": -1539.0,
					"date_acquired": "2018-02-02T16:19:15.693Z",
					"id": 284952,
					"quantity": -9.0,
					"symbol": "SPY180316P00266000"
				}, {
					"cost_basis": -1270.0,
					"date_acquired": "2018-02-06T17:39:07.118Z",
					"id": 286943,
					"quantity": -2.0,
					"symbol": "VXX180223C00050000"
				}, {
					"cost_basis": 982.0,
					"date_acquired": "2018-02-06T17:39:07.118Z",
					"id": 286944,
					"quantity": 2.0,
					"symbol": "VXX180223C00055000"
				}, {
					"cost_basis": 1290.0,
					"date_acquired": "2018-02-06T17:36:58.677Z",
					"id": 286939,
					"quantity": 2.0,
					"symbol": "VXX180302P00046000"
				}]
			}
		}]
	}
}
*/
