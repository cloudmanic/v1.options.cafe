//
// This is a websocket connection to the backend app. 
// Other than quotes all communication runs over this connection 

import { EventEmitter } from '@angular/core';
import { environment } from '../../../environments/environment';

declare var ClientJS: any;

export class AppService  
{  
  deviceId = ""  
  
  // Websocket Stuff
  ws = null;
  heartbeat = null;
  missed_heartbeats = 0;  
  
  // Emitters - Pushers
  wsReconnecting = new EventEmitter<boolean>();  
  
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
    
    console.log("HERE");
  }

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
      console.log(e.data);
      
      let msg = JSON.parse(e.data);
      
      // Is this a pong to our ping or some other return.
      if(msg.type == 'pong')
      {
        this.missed_heartbeats--;
      } else
      {        
        let msg_data = JSON.parse(msg.data);

        console.log(msg_data);
        
/*
        // Send quote to angular component
        switch(msg.type)
        {
          // User Profile refresh
          case 'UserProfile:refresh':
            this.doUserProfileRefresh(msg_data);     
          break;
          
          // Market Status refresh
          case 'MarketStatus:refresh':
            this.doMarketStatusRefresh(msg_data);     
          break;

          // Watchlist refresh
          case 'Watchlist:refresh':
            this.doWatchListRefresh(msg_data);     
          break;

          // Order refresh
          case 'Orders:refresh':
            this.doOrdersRefresh(msg_data);     
          break;
          
          // Balances refresh
          case 'Balances:refresh':
            this.doBalancesRefresh(msg_data);
          break;
        }
*/
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
              throw new Error('Too many missed heartbeats.');
            }
            
            this.ws.send(JSON.stringify({ type: 'ping' }));
            
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