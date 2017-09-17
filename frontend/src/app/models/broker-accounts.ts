export class BrokerAccounts {
  constructor(
    public AccountNumber: string,
    public Classification: string,
    public DayTrader: boolean,
    public OptionLevel: number,
    public Status: string,
    public Type: string     
  ){}
}
