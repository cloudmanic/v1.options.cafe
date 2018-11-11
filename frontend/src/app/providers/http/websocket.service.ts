//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
//
// This is a websocket connection to the backend app.
//

import { EventEmitter, Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { Order } from '../../models/order';
import { Balance } from '../../models/balance';
import { OrderLeg } from '../../models/order-leg';
import { ChangeDetected } from '../../models/change-detected';
import { MarketQuote } from '../../models/market-quote';
import { MarketStatus } from '../../models/market-status';
import { StateService } from '../state/state.service';

declare var ClientJS: any;

@Injectable()
export class WebsocketService  
{  
  deviceId = "";

  // Websocket Stuff
  ws = null;
  heartbeat = null;
  missed_heartbeats = 0;  
  
  // Emitters - Pushers
  wsReconnecting = new EventEmitter<boolean>();
  ordersPush = new EventEmitter<Order[]>();
  balancesPush = new EventEmitter<Balance[]>();
  marketStatusPush = new EventEmitter<MarketStatus>();
  changedDetectedPush = new EventEmitter<ChangeDetected>();
  quotePushData = new EventEmitter<MarketQuote>();    

  //
  // Construct!!
  //
  constructor(private stateService: StateService) 
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
    // Send quote to angular component
    switch(msg.uri)
    {  
      // Quote refresh
      case 'quote':
        this.doQuote(msg.body);
      break;

      // Orders refresh
      case 'orders':
        this.ordersPush.emit(Order.buildForEmit(msg.body));
      break;

      // Balances refresh
      case 'balances':
        this.balancesPush.emit(Balance.buildForEmit(msg.body));
      break;      

      // Market status
      case 'market-status':
        this.marketStatusPush.emit(MarketStatus.buildForEmit(msg.body));              
      break;

      // Change detected
      case 'change-detected':
        this.changedDetectedPush.emit(new ChangeDetected(msg.body.type));              
      break;              
    }
  }

  //
  // Manage the income quote.
  // 
  // {"type":"index","symbol":"$DJI","size":0,"last":25709.27,"open":25403.35,"high":25732.8,"low":25398.56,"bid":25717.47,"ask":25784.6,"close":25709.27,"prev_close":0,"change":399.28,"change_percentage":1.58,"volume":473357480,"average_volume":0,"last_volume":0,"description":"Dow Jones Industrial Average"}
  //
  private doQuote(msg_data: any)
  {
    // Build quote
    let quote = new MarketQuote(msg_data.last, msg_data.open, msg_data.ask, msg_data.bid, msg_data.prev_close, msg_data.symbol, msg_data.description, msg_data.change, msg_data.change_percentage);

    // Send quote out to components 
    this.quotePushData.emit(quote);

    // Store quote in cache
    this.stateService.SetQuote(quote);    
  }


  // ---------------------- Websocket Stuff ----------------------- //


  //
  // Setup normal data websocket connection.
  //
  setupWebSocket() 
  {
    // Setup websocket
    this.ws = new WebSocket(environment.ws_server + '/ws');
    
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