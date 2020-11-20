//
// Date: 2019-03-06
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2019-03-06
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package pull

import (
	"testing"

	"app.options.cafe/models"
	"github.com/nbio/st"
)

//
// endBookstrapping - Make sure we do not disable the broker.
//
func TestEndBookstrapping01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Users
	db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	db.Create(&models.User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	db.Create(&models.User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})

	// Test the before user
	user1, err := db.GetUserById(2)
	st.Expect(t, err, nil)
	st.Expect(t, user1.Bootstrapped, "No")

	// Test the function
	endBookstrapping(db, 2, "spicer+janewells@options.cafe")

	// Test the after user
	user2, err := db.GetUserById(2)
	st.Expect(t, err, nil)
	st.Expect(t, user2.Bootstrapped, "Yes")

	// Test the other user
	user3, err := db.GetUserById(3)
	st.Expect(t, err, nil)
	st.Expect(t, user3.Bootstrapped, "No")
}

//
// TestDoGetAllOrders01 - TODO(spicer): Make this work
//
func TestDoGetAllOrders01(t *testing.T) {
	// // Load config file.
	// env.ReadEnv("../../.env")
	//
	// // Start the db connection.
	// db, dbName, _ := models.NewTestDB("")
	// defer models.TestingTearDown(db, dbName)
	//
	// // Flush pending mocks after test execution
	// defer gock.Off()
	//
	// // Setup mock request.
	// gock.New("https://api.tradier.com/v1").
	// 	Get("/user/profile").
	// 	Reply(200).
	// 	BodyString(`{"profile":{"account":[{"account_number":"6YA05782","classification":"individual","date_created":"2016-08-01T21:09:15.000Z","day_trader":false,"option_level":0,"status":"active","type":"cash","last_update_date":"2016-08-01T21:09:15.000Z"},{"account_number":"6YA06984","classification":"traditional_ira","date_created":"2016-09-06T13:55:19.000Z","day_trader":false,"option_level":2,"status":"active","type":"cash","last_update_date":"2016-09-06T13:55:57.000Z"},{"account_number":"6YA06085","classification":"individual","date_created":"2016-08-01T21:09:15.000Z","day_trader":false,"option_level":4,"status":"active","type":"margin","last_update_date":"2016-08-01T21:09:15.000Z"}],"id":"id-2mxb7mwa","name":"Spicer Matthews"}}`)
	//
	// gock.New("https://api.tradier.com/v1").
	// 	Get("/accounts/6YA05782/orders").
	// 	Reply(200).
	// 	BodyString(`{"orders":{"order":[{"id":122758,"type":"market","symbol":"USO","side":"buy","quantity":20.00000000,"status":"filled","duration":"pre","avg_fill_price":8.47990000,"exec_quantity":20.00000000,"last_fill_price":8.47990000,"last_fill_quantity":20.00000000,"remaining_quantity":0.00000000,"create_date":"2016-02-23T17:01:13.993Z","transaction_date":"2016-02-23T17:01:14.000Z","class":"equity"},{"id":190500,"type":"market","symbol":"NFLX","side":"buy","quantity":4.00000000,"status":"filled","duration":"pre","avg_fill_price":98.94990000,"exec_quantity":4.00000000,"last_fill_price":98.94990000,"last_fill_quantity":4.00000000,"remaining_quantity":0.00000000,"create_date":"2016-09-16T17:17:34.855Z","transaction_date":"2016-09-16T17:17:35.490Z","class":"equity"},{"id":211123,"type":"limit","symbol":"DIS","side":"buy","quantity":4.00000000,"status":"filled","duration":"pre","price":96.83,"avg_fill_price":96.81990000,"exec_quantity":4.00000000,"last_fill_price":96.81990000,"last_fill_quantity":4.00000000,"remaining_quantity":0.00000000,"create_date":"2016-11-11T16:36:18.439Z","transaction_date":"2016-11-11T16:36:18.638Z","class":"equity"},{"id":245647,"type":"limit","symbol":"BAC","side":"buy","quantity":1.00000000,"status":"canceled","duration":"pre","price":20.0,"avg_fill_price":0.00000000,"exec_quantity":0.00000000,"last_fill_price":0.00000000,"last_fill_quantity":0.00000000,"remaining_quantity":0.00000000,"create_date":"2017-01-21T06:59:56.293Z","transaction_date":"2017-01-21T08:51:06.341Z","class":"equity"},{"id":245648,"type":"limit","symbol":"USO","side":"buy","quantity":1.00000000,"status":"canceled","duration":"gtc","price":11.2,"avg_fill_price":0.00000000,"exec_quantity":0.00000000,"last_fill_price":0.00000000,"last_fill_quantity":0.00000000,"remaining_quantity":0.00000000,"create_date":"2017-01-21T08:53:20.122Z","transaction_date":"2017-01-21T08:53:45.293Z","class":"equity"},{"id":255878,"type":"limit","symbol":"S","side":"buy","quantity":10.00000000,"status":"canceled","duration":"gtc","price":4.0,"avg_fill_price":0.00000000,"exec_quantity":0.00000000,"last_fill_price":0.00000000,"last_fill_quantity":0.00000000,"remaining_quantity":0.00000000,"create_date":"2017-02-05T06:09:37.501Z","transaction_date":"2017-03-07T20:27:33.969Z","class":"equity"},{"id":312418,"type":"limit","symbol":"VTI","side":"buy","quantity":4.00000000,"status":"filled","duration":"pre","price":121.21,"avg_fill_price":121.20990000,"exec_quantity":4.00000000,"last_fill_price":121.20990000,"last_fill_quantity":4.00000000,"remaining_quantity":0.00000000,"create_date":"2017-04-06T17:53:15.217Z","transaction_date":"2017-04-06T17:53:15.135Z","class":"equity"},{"id":313623,"type":"limit","symbol":"CONE","side":"buy","quantity":10.00000000,"status":"filled","duration":"pre","price":51.96,"avg_fill_price":51.95000000,"exec_quantity":10.00000000,"last_fill_price":51.95000000,"last_fill_quantity":10.00000000,"remaining_quantity":0.00000000,"create_date":"2017-04-07T17:55:25.603Z","transaction_date":"2017-04-07T17:55:58.279Z","class":"equity"}]}}`)
	//
	// gock.New("https://api.tradier.com/v1").
	// 	Get("/accounts/6YA06984/orders").
	// 	Reply(200).
	// 	BodyString(`{"orders":{"order":[{"id":122758,"type":"market","symbol":"USO","side":"buy","quantity":20.00000000,"status":"filled","duration":"pre","avg_fill_price":8.47990000,"exec_quantity":20.00000000,"last_fill_price":8.47990000,"last_fill_quantity":20.00000000,"remaining_quantity":0.00000000,"create_date":"2016-02-23T17:01:13.993Z","transaction_date":"2016-02-23T17:01:14.000Z","class":"equity"},{"id":190500,"type":"market","symbol":"NFLX","side":"buy","quantity":4.00000000,"status":"filled","duration":"pre","avg_fill_price":98.94990000,"exec_quantity":4.00000000,"last_fill_price":98.94990000,"last_fill_quantity":4.00000000,"remaining_quantity":0.00000000,"create_date":"2016-09-16T17:17:34.855Z","transaction_date":"2016-09-16T17:17:35.490Z","class":"equity"},{"id":211123,"type":"limit","symbol":"DIS","side":"buy","quantity":4.00000000,"status":"filled","duration":"pre","price":96.83,"avg_fill_price":96.81990000,"exec_quantity":4.00000000,"last_fill_price":96.81990000,"last_fill_quantity":4.00000000,"remaining_quantity":0.00000000,"create_date":"2016-11-11T16:36:18.439Z","transaction_date":"2016-11-11T16:36:18.638Z","class":"equity"},{"id":245647,"type":"limit","symbol":"BAC","side":"buy","quantity":1.00000000,"status":"canceled","duration":"pre","price":20.0,"avg_fill_price":0.00000000,"exec_quantity":0.00000000,"last_fill_price":0.00000000,"last_fill_quantity":0.00000000,"remaining_quantity":0.00000000,"create_date":"2017-01-21T06:59:56.293Z","transaction_date":"2017-01-21T08:51:06.341Z","class":"equity"},{"id":245648,"type":"limit","symbol":"USO","side":"buy","quantity":1.00000000,"status":"canceled","duration":"gtc","price":11.2,"avg_fill_price":0.00000000,"exec_quantity":0.00000000,"last_fill_price":0.00000000,"last_fill_quantity":0.00000000,"remaining_quantity":0.00000000,"create_date":"2017-01-21T08:53:20.122Z","transaction_date":"2017-01-21T08:53:45.293Z","class":"equity"},{"id":255878,"type":"limit","symbol":"S","side":"buy","quantity":10.00000000,"status":"canceled","duration":"gtc","price":4.0,"avg_fill_price":0.00000000,"exec_quantity":0.00000000,"last_fill_price":0.00000000,"last_fill_quantity":0.00000000,"remaining_quantity":0.00000000,"create_date":"2017-02-05T06:09:37.501Z","transaction_date":"2017-03-07T20:27:33.969Z","class":"equity"},{"id":312418,"type":"limit","symbol":"VTI","side":"buy","quantity":4.00000000,"status":"filled","duration":"pre","price":121.21,"avg_fill_price":121.20990000,"exec_quantity":4.00000000,"last_fill_price":121.20990000,"last_fill_quantity":4.00000000,"remaining_quantity":0.00000000,"create_date":"2017-04-06T17:53:15.217Z","transaction_date":"2017-04-06T17:53:15.135Z","class":"equity"},{"id":313623,"type":"limit","symbol":"CONE","side":"buy","quantity":10.00000000,"status":"filled","duration":"pre","price":51.96,"avg_fill_price":51.95000000,"exec_quantity":10.00000000,"last_fill_price":51.95000000,"last_fill_quantity":10.00000000,"remaining_quantity":0.00000000,"create_date":"2017-04-07T17:55:25.603Z","transaction_date":"2017-04-07T17:55:58.279Z","class":"equity"}]}}`)
	//
	// gock.New("https://api.tradier.com/v1").
	// 	Get("/accounts/6YA06085/orders").
	// 	Reply(200).
	// 	BodyString(`{"orders":{"order":[{"id":122758,"type":"market","symbol":"USO","side":"buy","quantity":20.00000000,"status":"filled","duration":"pre","avg_fill_price":8.47990000,"exec_quantity":20.00000000,"last_fill_price":8.47990000,"last_fill_quantity":20.00000000,"remaining_quantity":0.00000000,"create_date":"2016-02-23T17:01:13.993Z","transaction_date":"2016-02-23T17:01:14.000Z","class":"equity"},{"id":190500,"type":"market","symbol":"NFLX","side":"buy","quantity":4.00000000,"status":"filled","duration":"pre","avg_fill_price":98.94990000,"exec_quantity":4.00000000,"last_fill_price":98.94990000,"last_fill_quantity":4.00000000,"remaining_quantity":0.00000000,"create_date":"2016-09-16T17:17:34.855Z","transaction_date":"2016-09-16T17:17:35.490Z","class":"equity"},{"id":211123,"type":"limit","symbol":"DIS","side":"buy","quantity":4.00000000,"status":"filled","duration":"pre","price":96.83,"avg_fill_price":96.81990000,"exec_quantity":4.00000000,"last_fill_price":96.81990000,"last_fill_quantity":4.00000000,"remaining_quantity":0.00000000,"create_date":"2016-11-11T16:36:18.439Z","transaction_date":"2016-11-11T16:36:18.638Z","class":"equity"},{"id":245647,"type":"limit","symbol":"BAC","side":"buy","quantity":1.00000000,"status":"canceled","duration":"pre","price":20.0,"avg_fill_price":0.00000000,"exec_quantity":0.00000000,"last_fill_price":0.00000000,"last_fill_quantity":0.00000000,"remaining_quantity":0.00000000,"create_date":"2017-01-21T06:59:56.293Z","transaction_date":"2017-01-21T08:51:06.341Z","class":"equity"},{"id":245648,"type":"limit","symbol":"USO","side":"buy","quantity":1.00000000,"status":"canceled","duration":"gtc","price":11.2,"avg_fill_price":0.00000000,"exec_quantity":0.00000000,"last_fill_price":0.00000000,"last_fill_quantity":0.00000000,"remaining_quantity":0.00000000,"create_date":"2017-01-21T08:53:20.122Z","transaction_date":"2017-01-21T08:53:45.293Z","class":"equity"},{"id":255878,"type":"limit","symbol":"S","side":"buy","quantity":10.00000000,"status":"canceled","duration":"gtc","price":4.0,"avg_fill_price":0.00000000,"exec_quantity":0.00000000,"last_fill_price":0.00000000,"last_fill_quantity":0.00000000,"remaining_quantity":0.00000000,"create_date":"2017-02-05T06:09:37.501Z","transaction_date":"2017-03-07T20:27:33.969Z","class":"equity"},{"id":312418,"type":"limit","symbol":"VTI","side":"buy","quantity":4.00000000,"status":"filled","duration":"pre","price":121.21,"avg_fill_price":121.20990000,"exec_quantity":4.00000000,"last_fill_price":121.20990000,"last_fill_quantity":4.00000000,"remaining_quantity":0.00000000,"create_date":"2017-04-06T17:53:15.217Z","transaction_date":"2017-04-06T17:53:15.135Z","class":"equity"},{"id":313623,"type":"limit","symbol":"CONE","side":"buy","quantity":10.00000000,"status":"filled","duration":"pre","price":51.96,"avg_fill_price":51.95000000,"exec_quantity":10.00000000,"last_fill_price":51.95000000,"last_fill_quantity":10.00000000,"remaining_quantity":0.00000000,"create_date":"2017-04-07T17:55:25.603Z","transaction_date":"2017-04-07T17:55:58.279Z","class":"equity"}]}}`)
	//
	// // Shared vars we use.
	// ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)
	//
	// // Users
	// db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	// db.Create(&models.User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	// db.Create(&models.User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})
	//
	// // Brokers
	// db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts})
	// db.Create(&models.Broker{Name: "Tradeking", UserId: 1, AccessToken: "456", RefreshToken: "xyz", TokenExpirationDate: ts})
	// db.Create(&models.Broker{Name: "Etrade", UserId: 1, AccessToken: "789", RefreshToken: "mno", TokenExpirationDate: ts})
	//
	// // Get user.
	// user, err := db.GetUserById(1)
	// st.Expect(t, err, nil)
	//
	// // Get broker
	// broker, err := db.GetBrokerById(1)
	// st.Expect(t, err, nil)
	//
	// // Setup new API
	// api := &tradier.Api{}
	//
	// // Run get all orders
	// err = DoGetAllOrders(db, api, user, broker)
	// st.Expect(t, err, nil)

}

/* End File */
