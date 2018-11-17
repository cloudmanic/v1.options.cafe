//
// Date: 11/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Injectable, EventEmitter } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { AnalyzeResult, AnalyzeLeg } from '../../models/analyze-result';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class AnalyzeService  
{ 
  dialog = new EventEmitter<AnalyzeTrade>();

  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get a result by option and price 
  //
  getOptionsUnderlyingPriceResult(open: number, legs: AnalyzeLeg[]): Observable<AnalyzeResult[]> 
  {
    // Build post request
    let post = {
      open_cost: open,
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

export class AnalyzeTrade
{
  OpenCost: number;
  Legs: AnalyzeLeg[];
}

/* End File */