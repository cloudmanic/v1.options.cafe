export class Balance {
  constructor(
    public AccountNumber: string,
    public AccountValue: number,
    public TotalCash: number, 
    public OptionBuyingPower: number, 
    public StockBuyingPower: number, 
  ){}
}
