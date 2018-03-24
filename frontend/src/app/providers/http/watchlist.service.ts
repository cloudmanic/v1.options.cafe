//
// Date: 2/13/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Symbol } from '../../models/symbol';
import { Watchlist } from '../../models/watchlist';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable()
export class WatchlistService  
{ 
  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // Get watchlists
  //
  get() : Observable<Watchlist[]> {
    return this.http.get<Watchlist[]>(environment.app_server + '/api/v1/watchlists').map(
      (data) => { 
        
        let watchlists: Watchlist[] = []

        // Build data
        for(let i = 0; i < data.length; i++)
        {
          watchlists.push(Watchlist.buildForEmit(data[i])) 
        }

        return watchlists; 
      }
    );
  }

  //
  // Get watchlist by Id
  //
  getById(id: number) : Observable<Watchlist> {
    return this.http.get<Watchlist>(environment.app_server + '/api/v1/watchlists/' + id).map(
      (data) => { return Watchlist.buildForEmit(data); 
    });
  } 

  //
  // Add symbol to a watchlist by Id
  //
  addSymbolByWatchlistId(id: number, symbolId: number): Observable<Symbol> {

    let post = {
      symbol_id: symbolId
    }

    return this.http.post<Symbol>(environment.app_server + '/api/v1/watchlists/' + id + '/symbol', post)
      .map((data) => { return new Symbol().fromJson(data.symbol); });
  }    
}

/* End File */