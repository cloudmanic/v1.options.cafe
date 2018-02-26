//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

export class HistoricalQuote {
  
  public Date: Date;
  public Open: number;
  public Close: number;  
  public High: number;
  public Low: number;
  public Volume: number;

  //
  // Constructor...
  //
  constructor(
    Date: Date,
    Open: number,
    Close: number,
    High: number,
    Low: number,
    Volume: number
  ){
    this.Date = Date;
    this.Open = Open;
    this.Close = Close;    
    this.High = High;
    this.Low = Low;
    this.Volume = Volume;
  }

  //
  // Build object for emitting to the app.
  //
  public static buildForEmit(data) : HistoricalQuote[]  {

    let result = [];
    let date: Date;

    if(! data)
    {
      return result;      
    }

    for(let i = 0; i < data.length; i++)
    {
      if(data[i].time != "0001-01-01T00:00:00Z")
      {
        date = new Date(data[i].time);
      } else 
      {
        date = new Date(data[i].date + " 00:00:00");
      }

      result.push(new HistoricalQuote(
        date, 
        data[i].open,
        data[i].close,
        data[i].high, 
        data[i].low,
        data[i].volume
      ));
    }


    return result; 
  }  
}

/* End File */