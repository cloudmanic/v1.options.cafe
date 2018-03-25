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
  get() : Observable<Watchlist[]> 
  {
    return this.http.get<Watchlist[]>(environment.app_server + '/api/v1/watchlists').map(
      (data) => { 
        
        let watchlists: Watchlist[] = []

        // Build data
        for(let i = 0; i < data.length; i++)
        {
          watchlists.push(new Watchlist().fromJson(data[i])); 
        }

        return watchlists; 
      }
    );
  }

  //
  // Get watchlist by Id
  //
  getById(id: number) : Observable<Watchlist> 
  {
    return this.http.get<Watchlist>(environment.app_server + '/api/v1/watchlists/' + id)
      .map((data) => { return new Watchlist().fromJson(data); });
  } 

  //
  // Create a new watchlist
  //
  create(name: string): Observable<Watchlist> {
    let body = {
      name: name
    }

    return this.http.post<Watchlist>(environment.app_server + '/api/v1/watchlists', body)
      .map((data) => { return new Watchlist().fromJson(data); });
  }

  //
  // Update a watchlist by Id
  //
  update(id: number, name: string): Observable<boolean> 
  {
    let body = {
      name: name
    }

    return this.http.put(environment.app_server + '/api/v1/watchlists/' + id, body)
      .map((data) => { return true; });
  }

  //
  // Delete watchlist by Id
  //
  delete(id: number): Observable<boolean> {
    return this.http.delete(environment.app_server + '/api/v1/watchlists/' + id)
      .map((data) => { return true; });
  }   

  //
  // Add symbol to a watchlist by Id
  //
  addSymbolByWatchlistId(id: number, symbolId: number): Observable<Symbol> 
  {
    let post = {
      symbol_id: symbolId
    }

    return this.http.post<AddSymbolResponse>(environment.app_server + '/api/v1/watchlists/' + id + '/symbol', post)
      .map((data) => { return new Symbol().fromJson(data.symbol); });
  }

  //
  // Reorder the symbols in a watchlist by Id
  //
  reorder(id: number, ids: number[]): Observable<boolean> {
    let body = {
      ids: ids
    }

    return this.http.put(environment.app_server + '/api/v1/watchlists/' + id + '/reorder', body)
      .map((data) => { return true; });
  } 


  //
  // Delete a symbol from a watchlist
  //
  deleteSymbol(id: number, symbol: number): Observable<boolean> {
    return this.http.delete(environment.app_server + '/api/v1/watchlists/' + id + '/symbol/' + symbol)
      .map((data) => { return true; });
  }    
}

//
// Response of a add Symbol
//
interface AddSymbolResponse 
{
  symbol: Symbol
}

/* End File */