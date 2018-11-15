//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Symbol } from './symbol';
import { Position } from './position';

//
// TradeGroup Model
//
export class TradeGroup 
{
  public Id: number;
  public Name: string;
  public Status: string;
  public Type: string;
  public OpenDate: string;
  public ClosedDate: string;
  public Credit: number;
  public Profit: number;
  public Proceeds: number;  
  public PercentGain: number;
  public Risked: number;
  public Commission: number;
  public Note: string;
  public Positions: Position[];

  //
  // Constructor
  //
  constructor(
    Id: number, 
    Name: string,
    Status: string, 
    Type: string,
    OpenDate: string,
    ClosedDate: string,
    Credit: number,
    Profit: number,
    Proceeds: number,    
    PercentGain: number,
    Risked: number,
    Commission: number,
    Note: string,
    Positions: Position[]        
  ) {
    this.Id = Id;
    this.Name = Name;
    this.Status = Status;
    this.Type = Type;
    this.OpenDate = OpenDate;
    this.ClosedDate = ClosedDate;
    this.Credit = Credit;
    this.Profit = Profit;
    this.Proceeds = Proceeds;    
    this.PercentGain = PercentGain;
    this.Risked = Risked;
    this.Commission = Commission;
    this.Note = Note;
    this.Positions = Positions;
  }

  //
  // Build object for emitting to the app.
  //
  public static buildForEmit(data) : TradeGroup[]  
  {

    let tg = [];

    if(! data)
    {
      return tg;      
    }

    for(let i = 0; i < data.length; i++)
    {
      // Add in the positions
      let positions = [];
      
      for(let k = 0; k < data[i].positions.length; k++)
      {
        positions.push(new Position(
          data[i].positions[k].id,
          data[i].positions[k].open_date, 
          data[i].positions[k].close_date, 
          data[i].positions[k].qty,
          data[i].positions[k].org_qty,
          data[i].positions[k].cost_basis,
          data[i].positions[k].proceeds,
          data[i].positions[k].profit,
          new Symbol().fromJson(data[i].positions[k].symbol)                           
        ));
      }

      tg.push(new TradeGroup(
        data[i].id, 
        data[i].name,         
        data[i].status, 
        data[i].type, 
        data[i].open_date, 
        data[i].closed_date,
        data[i].credit,
        data[i].profit,
        data[i].proceeds,
        data[i].percent_gain,
        data[i].risked,
        data[i].commission,
        data[i].note,
        positions
       ));
    }

    return tg; 
  }
}

//
// Setup a class to hold all the different position types
//
export class TradeGroupsCont 
{
  public Equity: TradeGroup[] = []; 
  public ShortEquity: TradeGroup[] = []; 
  public Option: TradeGroup[] = []; 
  public PutCreditSpread: TradeGroup[] = []; 
  public CallCreditSpread: TradeGroup[] = []; 
  public PutDebitSpread: TradeGroup[] = []; 
  public CallDebitSpread: TradeGroup[] = []; 
  public IronCondor: TradeGroup[] = []; 
  public ReverseIronCondor: TradeGroup[] = [];
  public LongCallButterfly: TradeGroup[] = [];
  public LongPutButterfly: TradeGroup[] = [];   
  public Other: TradeGroup[] = []; 
}

/* End Find */