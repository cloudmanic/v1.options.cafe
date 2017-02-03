import { BrokerAccounts } from './broker-accounts';

export class UserProfile {
  constructor(
    public Id: string,
    public Name: string,
    public Accounts: BrokerAccounts[] 
  ){}
}
