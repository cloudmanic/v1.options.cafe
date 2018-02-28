export class MarketQuote {
  
  //
  // Constructor...
  //
  constructor(
    public last: number,
    public open: number,
    public ask: number,
    public bid: number,    
    public prev_close: number,
    public symbol: string,
    public description: string,
    public change: number,
    public change_percentage: number,
  ){}
  
  //
  // Return the class color
  //
  classColor() {
        
    if(this.change_percentage > 0)
    {
      return 'green';
    } else if(this.change_percentage < 0)
    {
      return 'red';
    }
    
    return '';
    
  }
}

/* End File */