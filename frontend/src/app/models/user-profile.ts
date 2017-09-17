//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { BrokerAccounts } from './broker-accounts';

export class UserProfile 
{
  
  //
  // Construct...
  //
  constructor(
    public Id: string,
    public Name: string,
    public Accounts: BrokerAccounts[] 
  ){}
  
  //
  // Build user profile
  //
  build(data) {
    
    // Clear accounts array.
    this.Accounts = [];    
    
    // Setup the array of accounts.
    for(var i in data.Accounts)
    {
      this.Accounts.push(new BrokerAccounts(
        data.Accounts[i].AccountNumber,
        data.Accounts[i].Classification,
        data.Accounts[i].DayTrader,
        data.Accounts[i].OptionLevel,
        data.Accounts[i].Status,
        data.Accounts[i].Type       
      ));
    }
    
    this.Id = data.Id;
    this.Name = data.Name;
        
  }
  
}

/* End File */