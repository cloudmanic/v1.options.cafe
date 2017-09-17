//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
//
// This is a websocket connection to the backend app. 
// Other than quotes all communication runs over this connection 
//

import { EventEmitter, Injectable } from '@angular/core';
import { AppState } from '../app.state.service';
import { environment } from '../../../environments/environment';
import { Order } from '../../models/order';
import { Balance } from '../../models/balance';
import { OrderLeg } from '../../models/order-leg';
import { Watchlist } from '../../models/watchlist';
import { WatchlistItems } from '../../models/watchlist-items';
import { MarketStatus } from '../../models/market-status';
import { UserProfile } from '../../models/user-profile';
import { BrokerAccount } from '../../models/broker-account';

declare var ClientJS: any;

@Injectable()
export class AppService  
{  
  deviceId = ""
  activeAccount = ""    
   
  // Websocket Stuff
  ws = null;
  heartbeat = null;
  missed_heartbeats = 0;  
  
  // Emitters - Pushers
  wsReconnecting = new EventEmitter<boolean>();
  ordersPush = new EventEmitter<Order[]>();
  balancesPush = new EventEmitter<Balance[]>();
  userProfilePush = new EventEmitter<UserProfile>();
  marketStatusPush = new EventEmitter<MarketStatus>();
  watchlistPush = new EventEmitter<Watchlist>();
  activeAccountPush = new EventEmitter<BrokerAccount>();    
  
  //
  // Construct!!
  //
  constructor(private appState: AppState) 
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
  processIncomingData(msg)
  {
    let msg_data = JSON.parse(msg.data);

    // console.log(msg_data);
        
    // Send quote to angular component
    switch(msg.type)
    {
      // User Profile refresh
      case 'UserProfile:refresh':
        this.userProfilePush.emit(UserProfile.buildForEmit(msg_data));  
      break;
      
      // Balances refresh
      case 'Balances:refresh':
        this.balancesPush.emit(Balance.buildForEmit(msg_data));
      break;      

      // Market Status refresh
      case 'MarketStatus:refresh':
        this.marketStatusPush.emit(MarketStatus.buildForEmit(msg_data));   
      break;
    
/*
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
*/
    }
    
  }
   
  
/*  
  //
  // Do watchlist Refresh
  //
  doWatchListRefresh (data) {
    
    // We only care about the default watchlist
    if(data.Id != 'default')
    {
      return false;
    }
    
    let ws = new Watchlist(data.Id, data.Name, []);

    for(var i in data.Symbols)
    {
      ws.items.push(new WatchlistItems(data.Symbols[i].id, data.Symbols[i].symbol));
    }
    
    this.watchlistPush.emit(ws);
    
  }
  
  //
  // Do refresh orders
  //
  doOrdersRefresh (data) {
    
    var orders = [];
    
    for(var i = 0; i < data.length; i++)
    {
      // Add in the legs
      var legs = [];
      
      if(data[i].NumLegs > 0)
      {
        for(var k = 0; k < data[i].Legs.length; k++)
        {
          legs.push(new OrderLeg(
            data[i].Legs[k].Type,
            data[i].Legs[k].Symbol,
            data[i].Legs[k].OptionSymbol, 
            data[i].Legs[k].Side, 
            data[i].Legs[k].Quantity, 
            data[i].Legs[k].Status, 
            data[i].Legs[k].Duration, 
            data[i].Legs[k].AvgFillPrice, 
            data[i].Legs[k].ExecQuantity, 
            data[i].Legs[k].LastFillPrice, 
            data[i].Legs[k].LastFillQuantity, 
            data[i].Legs[k].RemainingQuantity, 
            data[i].Legs[k].CreateDate, 
            data[i].Legs[k].TransactionDate          
          ));
        }
      }
      
      // Push the order on
      orders.push(new Order(
          data[i].Id,
          data[i].AccountId,
          data[i].AvgFillPrice,
          data[i].Class,
          data[i].CreateDate,
          data[i].Duration,
          data[i].ExecQuantity,
          data[i].LastFillPrice,
          data[i].LastFillQuantity,
          data[i].NumLegs,
          data[i].Price,
          data[i].Quantity,
          data[i].RemainingQuantity,
          data[i].Side,
          data[i].Status,
          data[i].Symbol,
          data[i].TransactionDate,
          data[i].Type,
          legs));
               
    }
    
    this.ordersPush.emit(orders);
    
  }
*/  

  // ------------------------ Push Data Back To Backend --------------------- //
  
  //
  // Request the backend sends all data again. (often do this on state change or page change)
  //
  requestAllData() {
    this.ws.send(JSON.stringify({  type: 'refresh-all-data', data: {} }));   
  }
  
  //
  // Set the active account.
  //
  setActiveAccount(account) {
    this.activeAccount = account;
    this.appState.setActiveAccount(account);  
    this.activeAccountPush.emit(account);
    this.requestAllData();
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
      if(msg.type == 'pong')
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