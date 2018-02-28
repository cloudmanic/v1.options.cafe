//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

export class BrokerAccount {

  public Id: number;
  public Name: string;
  public BrokerId: number;
  public AccountNumber: string;
  public StockCommission: number;
  public StockMin: number;
  public OptionCommission: number;
  public OptionSingleMin: number;
  public OptionMultiLegMin: number;
  public OptionBase: number;  

  //
  // Construct.
  //
  constructor(
    Id: number, 
    Name: string,
    BrokerId: number, 
    AccountNumber: string,
    StockCommission: number,
    StockMin: number,
    OptionCommission: number,
    OptionSingleMin: number,
    OptionMultiLegMin: number,
    OptionBase: number    
  ){
    this.Id = Id;
    this.Name = Name;
    this.BrokerId = BrokerId;
    this.AccountNumber = AccountNumber;
    this.StockCommission = StockCommission;
    this.StockMin = StockMin;
    this.OptionCommission = OptionCommission;
    this.OptionSingleMin = OptionSingleMin;
    this.OptionMultiLegMin = OptionMultiLegMin;
    this.OptionBase = OptionBase;     
  }
}

/* End File */
