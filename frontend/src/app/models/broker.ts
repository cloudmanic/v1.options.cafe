//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { BrokerAccount } from './broker-account';

export class Broker {

  public Id: number;
  public Name: string;
  public Status: string;
  public DisplayName: string;
  public BrokerAccounts: BrokerAccount[];
  public SettingsActiveBrokerAccount: BrokerAccount;
  
  //
  // Construct.
  //
  constructor(
    Id: number,
    Name: string,
    DisplayName: string,    
    Status: string,   
    BrokerAccounts: BrokerAccount[],
  ){
    this.Id = Id;
    this.Name = Name;
    this.DisplayName = DisplayName;
    this.Status = Status;    
    this.BrokerAccounts = BrokerAccounts;
  }

  //
  // Json to Object.
  //
  public static fromJson(json: Object): Broker 
  {
    let o = new Broker(json["id"], json["name"], json["display_name"], json["status"], []);
    return o;
  }

  //
  // Build object for emitting to the app.
  //
  public static buildForEmit(data) : Broker[]  {

    let result = [];

    if(! data)
    {
      return result;      
    }

    for(let i = 0; i < data.length; i++)
    {
      // Add in the positions
      let accounts = [];
      
      for(let k = 0; k < data[i].broker_accounts.length; k++)
      {
        accounts.push(new BrokerAccount(
          data[i].broker_accounts[k].id,
          data[i].broker_accounts[k].name,
          data[i].broker_accounts[k].broker_id,  
          data[i].broker_accounts[k].account_number,           
          data[i].broker_accounts[k].stock_commission, 
          data[i].broker_accounts[k].stock_min,
          data[i].broker_accounts[k].option_commission,
          data[i].broker_accounts[k].option_single_min,
          data[i].broker_accounts[k].option_multi_leg_min,
          data[i].broker_accounts[k].option_base                         
        ));
      }

      result.push(new Broker(
        data[i].id,
        data[i].name,
        data[i].display_name,         
        data[i].status,       
        accounts
       ));
    }

    return result; 
  }  
}

/* End File */