//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

export class Position {
  constructor(
    public Id: number,
    public Symbol: string,
    public OpenDate: string, 
    public ClosedDate: string,
    public Qty: number,
    public OrgQty: number,
    public CostBasis: number
  ){}
}

/* End Find */

