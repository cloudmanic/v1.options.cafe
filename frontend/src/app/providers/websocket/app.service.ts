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
import { Watchlist } from '../../models/watchlist';
import { WatchlistItems } from '../../models/watchlist-items';
import { MarketStatus } from '../../models/market-status';
import { UserProfile } from '../../models/user-profile';
import { BrokerAccount } from '../../models/broker-account';

declare var ClientJS: any;

//@Injectable()
export class AppService  
{  
  public deviceId = ""
  public activeAccount: BrokerAccount;    
   
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
  userProfilePush = new EventEmitter<UserProfile>();
  marketStatusPush = new EventEmitter<MarketStatus>();
  watchlistPush = new EventEmitter<Watchlist>();
  activeAccountPush = new EventEmitter<BrokerAccount>();    
  
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
  
  // ---------------------- Geters / Setters ----------------------- //
  
  //
  // Set active account.
  //
  public setActiveAccount(account: BrokerAccount) {
    
    this.activeAccount = account;
    localStorage.setItem('active_account', account.AccountNumber);
    this.activeAccountPush.emit(account);
    this.requestAllData();
    
  }

  //
  // Get active account.
  //
  public getActiveAccount() : BrokerAccount {
    return this.activeAccount;
  }

  // ---------------------- Incoming Data  ------------------------- //

  //
  // Process incoming data.
  //
  private processIncomingData(msg)
  {
    let msg_data = JSON.parse(msg.data);

    // console.log(msg_data);
        
    // Send quote to angular component
    switch(msg.type)
    {
      // User Profile refresh
      case 'UserProfile:refresh':
        var obj = UserProfile.buildForEmit(msg_data);
        this.calcActiveAccount(obj);
        this.userProfilePush.emit(obj);  
      break;
      
      // Balances refresh
      case 'Balances:refresh':
        this.balancesPush.emit(Balance.buildForEmit(msg_data));
      break;      

      // Market Status refresh
      case 'MarketStatus:refresh':
        this.marketStatusPush.emit(MarketStatus.buildForEmit(msg_data));   
      break;
    
      // Order refresh
      case 'Orders:refresh':
        this.orders = Order.buildForEmit(msg_data);       
        this.ordersPush.emit(this.orders);              
      break;

      // Watchlist refresh
      case 'Watchlist:refresh':
        
        console.log(msg_data);
        
        //this.doWatchListRefresh(msg_data);     
      break;
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
*/  

  // ------------------------ Helper Functions ------------------------------ //

  //
  // Figure out what our active account is based on data we passed in.
  //
  private calcActiveAccount(user: UserProfile) {

    // If we already have an active account do nothing.
    if(this.getActiveAccount())
    {
      return;
    }

    if(! user.Accounts.length)
    {     
      return;
    }

    if((! localStorage.getItem('active_account')) && user.Accounts.length)
    {      
      this.setActiveAccount(user.Accounts[0]);
      return;
    }
    
    var acn = localStorage.getItem('active_account');
      
    for(var i = 0; i < user.Accounts.length; i++)
    {
      if(user.Accounts[i].AccountNumber == acn)
      {
        this.setActiveAccount(user.Accounts[i]);            
      }
    }
         
  }

  // ------------------------ Push Data Back To Backend --------------------- //
  
  //
  // Request the backend sends all data again. (often do this on state change or page change)
  //
  public requestAllData() {
    this.ws.send(JSON.stringify({  type: 'refresh-all-data', data: {} }));   
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