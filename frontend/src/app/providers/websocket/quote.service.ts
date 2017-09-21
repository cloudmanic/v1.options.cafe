//
// Date: 9/21/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
//
// This is a websocket connection to the backend app for quotes.
//

import { EventEmitter } from '@angular/core';
import { environment } from '../../../environments/environment';
import { MarketQuote } from '../../models/market-quote';

declare var ClientJS: any;

export class QuoteService {
  
  deviceId = ""
  quotes = {};
  
  // Websocket Stuff
  ws = null;
  heartbeat = null;
  missed_heartbeats = 0;
  
  // Emitters
  wsReconnecting = new EventEmitter<boolean>();
  marketQuotePushData = new EventEmitter<MarketQuote>();

  //
  // Construct!!
  //
  constructor() {
    
    // Set the device id
    var clientJs = new ClientJS();
    this.deviceId = clientJs.getFingerprint();
 
    // Setup standard websocket connection.
    this.setupWebSocket();    

  }

  // ------------------------ Websocket Stuff --------------------- //
  
  //
  // Setup normal data websocket connection.
  //
  setupWebSocket() {
    
    // Setup websocket
    this.ws = new WebSocket(environment.ws_server + '/ws/quotes');
    
    // Websocket sent data to us.
    this.ws.onmessage = (e) =>
    { 
      let msg = JSON.parse(e.data);
      
      // Is this a pong to our ping or some other return.
      if(msg.type == 'pong')
      {
        this.missed_heartbeats--;
      } else
      {        
        // Send quote to angular component
        switch(msg.type)
        {
          // Real-time market quote
          case 'trade':
            
            // Have we seen this quote before?
            if(typeof this.quotes[msg.symbol] == "undefined")
            {
              this.quotes[msg.symbol] = new MarketQuote(msg.last, 0, 0, msg.symbol, '');
            } else
            {
              this.quotes[msg.symbol].last = msg.last;             
            }

            this.marketQuotePushData.emit(this.quotes[msg.symbol]); 
                     
          break;
          
          // DetailedQuotes refresh
          case 'DetailedQuotes:refresh':
            
            let msg_data = JSON.parse(msg.data);
            
            // Have we seen this quote before?
            if(typeof this.quotes[msg_data.Symbol] == "undefined")
            {
              this.quotes[msg_data.Symbol] = new MarketQuote(msg_data.Last, msg_data.Open, msg_data.PrevClose, msg_data.Symbol, msg_data.Description);
            } else
            {
              this.quotes[msg_data.Symbol].last = msg_data.Last;              
              this.quotes[msg_data.Symbol].open = msg_data.Open;
              this.quotes[msg_data.Symbol].prev_close = msg_data.PrevClose;              
              this.quotes[msg_data.Symbol].description = msg_data.Description;              
            }            
              
            this.marketQuotePushData.emit(this.quotes[msg_data.Symbol]);
            
          break;          
        }
      }      
    }
    
    // On Websocket open
    this.ws.onopen = (e) =>
    {    
      // Send Access Token (Give a few moments to get started)
      setTimeout(() => { 
        this.ws.send(JSON.stringify({ 
          type: 'set-access-token', 
          data: { access_token: localStorage.getItem('access_token'), device_id: this.deviceId }
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
              throw new Error('Too many missed heartbeats (quotes).');
            }
            
            this.ws.send(JSON.stringify({ type: 'ping' }));
            
          } catch(e) 
          {
            this.wsReconnecting.emit(true);
            clearInterval(this.heartbeat);
            this.heartbeat = null;
            console.warn("Closing connection (quotes). Reason: " + e.message);
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
