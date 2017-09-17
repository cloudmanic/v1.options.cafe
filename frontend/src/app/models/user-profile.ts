//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { BrokerAccount } from './broker-account';

export class UserProfile 
{
  
  //
  // Construct...
  //
  constructor(
    public Id: string,
    public Name: string,
    public Accounts: BrokerAccount[] 
  ){}
  
  //
  // Build user profile for emitting to the app.
  //
  public static buildForEmit(data) : UserProfile  {
    
    let user = new UserProfile(data.Id, data.Name, []); 
       
    // Setup the array of accounts.
    for(let i in data.Accounts)
    {
      user.Accounts.push(new BrokerAccount(
        data.Accounts[i].AccountNumber,
        data.Accounts[i].Classification,
        data.Accounts[i].DayTrader,
        data.Accounts[i].OptionLevel,
        data.Accounts[i].Status,
        data.Accounts[i].Type       
      ));
    }
     
    return user;
  }
  
}

/* End File */