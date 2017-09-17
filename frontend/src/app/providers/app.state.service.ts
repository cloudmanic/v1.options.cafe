//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { UserProfile } from '../models/user-profile';
import { BrokerAccount } from '../models/broker-account';

export class AppState  
{ 
  public userProfile: UserProfile;
    
  public accounts: BrokerAccount[];
  public activeAccount: BrokerAccount;

  //
  // Set the user profile for this app.
  //
  public setUserProfile(profile: UserProfile) {    
    this.userProfile = profile;
  }

  //
  // Get the user profile for this app.
  //
  public getUserProfile() : UserProfile {    
    return this.userProfile;
  }
  
  //
  // Set the accounts for this app.
  //
  public setAccounts(accounts: BrokerAccount[]) {    
    this.accounts = accounts;
  }
  
  //
  // Set the active account.
  //
  public setActiveAccount(activeAccount: BrokerAccount) {    
    this.activeAccount = activeAccount;
  }
  
  //
  // Get the active account.
  //
  public getActiveAccount() : BrokerAccount {
    return this.activeAccount;
  }    
}

/* End File */