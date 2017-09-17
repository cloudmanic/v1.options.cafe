//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//


export class Balance {
  
  //
  // Construct.
  //
  constructor(
    public AccountNumber: string,
    public AccountValue: number,
    public TotalCash: number, 
    public OptionBuyingPower: number, 
    public StockBuyingPower: number, 
  ){}
  
  //
  // Build build the data for emitting to the app. 
  //
  public static buildForEmit(data) : Balance[] {
    
    let balances = [];
    
    for(let i = 0; i < data.length; i++)
    {
      balances.push(new Balance(
        data[i].AccountNumber,
        data[i].AccountValue,
        data[i].TotalCash,
        data[i].OptionBuyingPower,
        data[i].StockBuyingPower        
      ));               
    }
    
    return balances;
        
  }  
  
}

/* End File */
