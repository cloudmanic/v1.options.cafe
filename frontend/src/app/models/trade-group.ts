//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Position } from './position';

//
// TradeGroup Model
//
export class TradeGroup {
  public Id: number;
  public Status: string;
  public Type: string;
  public OpenDate: string;
  public ClosedDate: string;
  public Commission: number;
  public Note: string;
  public Positions: []Position: 

  //
  // Constructor
  //
  constructor(
    Id: number, 
    Status: string, 
    Type: string,
    OpenDate: string,
    ClosedDate: string,
    Commission: number,
    Note: string,
    Positions: []Position        
  ) {
    this.Id = Id;
    this.Status = Status;
    this.Type = Type;
    this.OpenDate = OpenDate;
    this.ClosedDate = ClosedDate;
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
          data[i].positions[k].symbol,
          data[i].positions[k].open_date, 
          data[i].positions[k].close_date, 
          data[i].positions[k].qty,
          data[i].positions[k].org_qty,
          data[i].positions[k].cost_basis                 
        ));
      }

      tg.push(new TradeGroup(
        data[i].id, 
        data[i].status, 
        data[i].type, 
        data[i].open_date, 
        data[i].closed_date,
        data[i].commission,
        data[i].note,
        positions
       ));
    }

    return tg; 
  }
}

/* End Find */