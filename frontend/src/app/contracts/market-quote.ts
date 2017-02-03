export class MarketQuote {
  
  //
  // Constructor...
  //
  constructor(
    public last: number,
    public open: number,
    public prev_close: number,
    public symbol: string,
    public description: string
  ){}

  //
  // Return the daily change
  //
  dailyChange() {
    return (this.last - this.prev_close).toFixed(2);
  }
  
  //
  // Return the percent change
  //
  percentChange() {
    
    // Do we have data yet?
    if(this.prev_close <= 0)
    {
      return 0;
    }
    
    return parseFloat((((this.last - this.prev_close) / this.prev_close) * 100).toFixed(2));
  }
  
  //
  // Return the class color
  //
  classColor() {
    
    let change = this.percentChange();
    
    if(change > 0)
    {
      return 'green';
    } else if(change < 0)
    {
      return 'red';
    }
    
    return '';
    
  }
}

/* End File */