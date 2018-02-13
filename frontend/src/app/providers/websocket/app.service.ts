//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
//
// This is a websocket connection to the backend app. 
// Other than quotes all communication runs over this connection 
//

import { EventEmitter /* Injectable */ } from '@angular/core';
import { environment } from '../../../environments/environment';
import { Order } from '../../models/order';
import { Balance } from '../../models/balance';
import { OrderLeg } from '../../models/order-leg';
import { MarketStatus } from '../../models/market-status';

declare var ClientJS: any;

//@Injectable()
export class AppService  
{  
  public deviceId = ""
  //public activeAccount: BrokerAccount;    
   
  // Cache some of the data.
  public orders: Order[];
   
  // Websocket Stuff
  ws = null;
  heartbeat = null;
  missed_heartbeats = 0;  
  
  // Emitters - Pushers
  wsReconnecting = new EventEmitter<boolean>();
  ordersPush = new EventEmitter<Order[]>();
  balancesPush = new EventEmitter<Balance[]>();
  marketStatusPush = new EventEmitter<MarketStatus>();
  
  //
  // Construct!!
  //
  constructor() 
  {
    // Set the device id
    var clientJs = new ClientJS();
    this.deviceId = clientJs.getFingerprint();
        
    // Setup standard websocket connection.
    this.setupWebSocket();
  }
  
  // ---------------------- Incoming Data  ------------------------- //

  //
  // Process incoming data.
  //
  private processIncomingData(msg)
  {
    let msg_data = JSON.parse(msg.body);

    // console.log(msg_data);
        
    // Send quote to angular component
    switch(msg.uri)
    {      
      // Balances refresh
      case 'balances':
        this.balancesPush.emit(Balance.buildForEmit(msg_data));
      break;      

      // Market Status refresh
      case 'market/status':
        this.marketStatusPush.emit(MarketStatus.buildForEmit(msg_data));   
      break;
    
      // Order refresh
      case 'orders':
        this.orders = Order.buildForEmit(msg_data);       
        this.ordersPush.emit(this.orders);              
      break;  
    }
    
  }

  // ------------------------ Push Data Back To Backend --------------------- //
  
  //
  // Request the backend sends all data again. (often do this on state change or page change)
  //
  public RequestAllData() {
    this.ws.send(JSON.stringify({ uri: 'data/all', body: {} }));   
  }
  
  //
  // Request the backend sends watchlist data.
  //
  public RequestWatchlistData() {
    this.ws.send(JSON.stringify({ uri: 'watchlists', body: {} }));   
  }
 
  // ---------------------- Websocket Stuff ----------------------- //


  //
  // Setup normal data websocket connection.
  //
  setupWebSocket() 
  {
    // Setup websocket
    this.ws = new WebSocket(environment.ws_server + '/ws/core');
    
    // Websocket sent data to us.
    this.ws.onmessage = (e) =>
    { 
      //console.log(e.data);
      
      let msg = JSON.parse(e.data);
      
      // Is this a pong to our ping or some other return.
      if(msg.uri == 'pong')
      {
        this.missed_heartbeats--;
      } else
      {        
        this.processIncomingData(msg);
      }
    }
    
    // On Websocket open
    this.ws.onopen = (e) =>
    {      
      // Send Access Token (Give a few moments to get started)
      setTimeout(() => { 

        // Send auth
        this.ws.send(JSON.stringify({ 
          uri: 'set-access-token', 
          body: { access_token: localStorage.getItem('access_token'), device_id: this.deviceId }
        }));

        // Request all data.
        this.RequestAllData()

      }, 1000);
      
      // Tell the UI we are connected
      this.wsReconnecting.emit(false);
      
      // Setup the connection heartbeat
      if(this.heartbeat === null) 
      {
        this.missed_heartbeats = 0;
        
        this.heartbeat = setInterval(() => {
         
          try {
            this.missed_heartbeats++;
            
            if(this.missed_heartbeats >= 5)
            {
              throw new Error('Too many missed heartbeats.');
            }
            
            this.ws.send(JSON.stringify({ uri: 'ping' }));
            
          } catch(e) 
          {
            this.wsReconnecting.emit(true);
            clearInterval(this.heartbeat);
            this.heartbeat = null;
            console.warn("Closing connection. Reason: " + e.message);
            this.ws.close();
          }
          
        }, 5000);
        
      } else
      {
        clearInterval(this.heartbeat);
      }      
    }
    
    // On Close
    this.ws.onclose = () => 
    {      
      // Kill Ping heartbeat.
      clearInterval(this.heartbeat);
      this.heartbeat = null;
      this.ws = null;
      
      // Try to reconnect
      this.wsReconnecting.emit(true);
      setTimeout(() => { this.setupWebSocket(); }, 3 * 1000);
    }
  }

}

/* End File */