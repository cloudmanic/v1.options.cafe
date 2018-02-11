//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Symbol } from './symbol';

export class Position {

  //
  // Construct.
  //
  constructor(
    public Id: number,
    public OpenDate: string, 
    public ClosedDate: string,
    public Qty: number,
    public OrgQty: number,
    public CostBasis: number,
    public Symbol: Symbol    
  ){}
}

/* End Find */

