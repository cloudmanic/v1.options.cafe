//
// Date: 11/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { AnalyzeResult, AnalyzeLeg } from '../../models/analyze-result';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class AnalyzeService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get a result by option and price 
  //
  getOptionsUnderlyingPriceResult(open: number, price: number, legs: AnalyzeLeg[]): Observable<AnalyzeResult[]> 
  {
    // Build post request
    let post = {
      open_cost: open,
      current_underlying_price: price,
      legs: [] 
    }

    // Add in legs
    for(let i = 0; i < legs.length; i++)
    {
      post.legs.push({ symbol_str: legs[i].SymbolStr, qty: legs[i].Qty });
    }

    return this.http.post<AnalyzeResult[]>(environment.app_server + '/api/v1/analyze/options/underlying-price', post)
      .map((data) => { return new AnalyzeResult().fromJsonList(data); });
  } 

}

/* End File */