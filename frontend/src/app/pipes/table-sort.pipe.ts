//
// Date: 2018-11-04
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-04
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'tableSort'
})

export class TableSortPipe implements PipeTransform 
{
  //
  // Transform
  //
  transform(results: any[], col: string, order: number): any[] 
  {

    // Check if is not null
    if (!results) 
    {
      return results;
    }

    return results.sort((a: any, b: any) => {

      // Order * (-1): We change our order
      return a[col] > b[col] ? order : order * (- 1);

    });
  }

}

/* End File */