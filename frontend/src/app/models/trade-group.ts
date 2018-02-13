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
export class TradeGroup {
  public Id: number;
  public Name: string;
  public Status: string;
  public Type: string;
  public OpenDate: string;
  public ClosedDate: string;
  public Gain: number, 
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
    Profit: number,
    Gain: number,
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
    this.Profit = Profit;
    this.Gain = Gain;
    this.Risked = Risked;
    this.Commission = Commission;
    this.Note = Note;
    this.Positions = Positions;
  }

  //
  // Build object for emitting to the app.
  //
  public static buildForEmit(data) : TradeGroup[]  {

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
          new Symbol(data[i].positions[k].symbol.id, data[i].positions[k].symbol.name, data[i].positions[k].symbol.short_name)                           
        ));
      }

      tg.push(new TradeGroup(
        data[i].id, 
        data[i].name,         
        data[i].status, 
        data[i].type, 
        data[i].open_date, 
        data[i].closed_date,
        data[i].profit,
        data[i].gain,
        data[i].risked,
        data[i].commission,
        data[i].note,
        positions
       ));
    }

    return tg; 
  }
}

/* End Find */