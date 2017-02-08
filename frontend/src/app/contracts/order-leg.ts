export class OrderLeg {
  constructor(
    public Type: string,
    public Symbol: string,
    public OptionSymbol: string, 
    public Side: string, 
    public Quantity: number, 
    public Status: string, 
    public Duration: string, 
    public AvgFillPrice: number, 
    public ExecQuantity: number, 
    public LastFillPrice: number, 
    public LastFillQuantity: number, 
    public RemainingQuantity: number, 
    public CreateDate: string, 
    public TransactionDate: string,
  ){}
}
