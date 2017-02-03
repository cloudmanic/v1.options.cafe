import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { BrokerService } from '../services/broker.service';
import { UserProfile } from '../contracts/user-profile';
import { BrokerAccounts } from '../contracts/broker-accounts';

@Component({
  selector: 'oc-accounts',
  templateUrl: './accounts.component.html'
})
export class AccountsComponent implements OnInit {
  userProfile: UserProfile;
  selectedAccount: BrokerAccounts;

  //
  // Constructor....
  //
  constructor(private broker: BrokerService, private changeDetect: ChangeDetectorRef) { }

  //
  // OnInit....
  //
  ngOnInit() {

    // Subscribe to data updates from the broker - Market Status
    this.broker.userProfilePushData.subscribe(data => {
      this.userProfile = data;
      
      // Do we have an account already? Always have to reset the selected one when we get new account data.
      if((! this.selectedAccount) && (this.userProfile.Accounts.length))
      {
        this.selectedAccount = this.userProfile.Accounts[0];
        this.broker.setActiveAccountId(this.selectedAccount.AccountNumber);
      } else
      {
        for(var i = 0; i < this.userProfile.Accounts.length; i++)
        {
          if(this.userProfile.Accounts[i].AccountNumber == this.selectedAccount.AccountNumber)
          {
            this.selectedAccount = this.userProfile.Accounts[i];            
          }
        }
      }
      
      this.changeDetect.detectChanges();
    });

  }
  
  //
  // On account change.
  //
  onAccountChange() {
    
    this.broker.setActiveAccountId(this.selectedAccount.AccountNumber);
  
  }

}
