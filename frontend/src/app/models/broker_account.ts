//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

export class BrokerAccount {
  
  //
  // Construct.
  //
  constructor(
    public Id: number,
    public Name: string,
    public BrokerAccounts: BrokerAccount[]
  ){}
}

/* End File */