export class Order {
  constructor(
    public Id: number,
    public AvgFillPrice: number,
    public Class: string,
    public CreateDate: string,
    public Duration: string,
    public ExecQuantity: string,
    public LastFillPrice: number,
    public LastFillQuantity: number,
    public NumLegs: number,
    public Price: number,
    public Quantity: number,
    public RemainingQuantity: number,
    public Side: string,
    public Status: string,
    public Symbol: string,
    public TransactionDate: string,
    public Type: string
  ){}
}
