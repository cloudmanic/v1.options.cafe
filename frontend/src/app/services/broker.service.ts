import { EventEmitter } from '@angular/core';
import { Order } from '../contracts/order';
import { Balance } from '../contracts/balance';
import { OrderLeg } from '../contracts/order-leg';
import { Watchlist } from '../contracts/watchlist';
import { WatchlistItems } from '../contracts/watchlist-items';
import { MarketStatus } from '../contracts/market-status';
import { UserProfile } from '../contracts/user-profile';
import { BrokerAccounts } from '../contracts/broker-accounts';

/*
declare var electron: any;
declare var tradier_api_key: any;
*/

declare var ws_server: any;
declare var Fingerprint2: any;

export class BrokerService {
  
  deviceId = ""
  activeAccount = ""
  
  // Data objects
  marketStatus = new MarketStatus('', '');
  userProfile = new UserProfile('', '', []);
  
  // Websocket Stuff
  ws = null;
  heartbeat = null;
  missed_heartbeats = 0;
  
  // Emitters - Pushers
  wsReconnecting = new EventEmitter<boolean>();
  ordersPushData = new EventEmitter<Order[]>();
  balancesPushData = new EventEmitter<Balance[]>();
  userProfilePushData = new EventEmitter<UserProfile>();
  marketStatusPushData = new EventEmitter<MarketStatus>();
  watchlistPushData = new EventEmitter<Watchlist>();
  activeAccountPushData = new EventEmitter<string>();   

  //
  // Construct!!
  //
  constructor() {
    
    var self = this;

    // Set the device id
    new Fingerprint2().get(function(result, components) {
      self.deviceId = result;
    });

    // Setup standard websocket connection.
    this.setupWebSocket();    

  }
  
  //
  // Do User Profile Refresh
  //
  doUserProfileRefresh (data) {
    
    // Clear accounts array.
    this.userProfile.Accounts = [];
    
    // Setup the array of accounts.
    for(var i in data.Accounts)
    {
      this.userProfile.Accounts.push(new BrokerAccounts(
        data.Accounts[i].AccountNumber,
        data.Accounts[i].Classification,
        data.Accounts[i].DayTrader,
        data.Accounts[i].OptionLevel,
        data.Accounts[i].Status,
        data.Accounts[i].Type       
      ));
    }
    
    this.userProfile.Id = data.Id;
    this.userProfile.Name = data.Name;
    this.userProfilePushData.emit(this.userProfile);
    
  }
  
  //
  // Do Market Status Refresh
  //
  doMarketStatusRefresh (data) {
    this.marketStatus.state = data.State;
    this.marketStatus.description = data.Description;
    this.marketStatusPushData.emit(this.marketStatus);
  }
  
  //
  // Do Balances Refresh
  //
  doBalancesRefresh (data) {
    
    var balances = [];
    
    for(var i = 0; i < data.length; i++)
    {
      balances.push(new Balance(
        data[i].AccountNumber,
        data[i].AccountValue,
        data[i].TotalCash,
        data[i].OptionBuyingPower,
        data[i].StockBuyingPower        
      ));               
    }  

    this.balancesPushData.emit(balances);
    
  }  
  
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
    
    this.watchlistPushData.emit(ws);
    
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
    
    this.ordersPushData.emit(orders);
    
  }

  // ------------------------ Push Data Back To Backend --------------------- //
  
  //
  // Set the active account id.
  //
  setActiveAccountId(account_id) {
    
    this.activeAccount = account_id;
    this.activeAccountPushData.emit(account_id);
  
  }

  // ------------------------ Websocket Stuff --------------------- //
  
  //
  // Setup normal data websocket connection.
  //
  setupWebSocket() {
    
    // Setup websocket
    this.ws = new WebSocket(ws_server + '/ws/core');
    
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
        let msg_data = JSON.parse(msg.data);
        
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
