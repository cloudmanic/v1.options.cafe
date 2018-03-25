//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { NgForm } from '@angular/forms';
import { Component, OnInit, ElementRef } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { SortablejsOptions } from 'angular-sortablejs';
import { Symbol } from '../../../models/symbol';
import { Watchlist } from '../../../models/watchlist';
import { WatchlistService } from '../../../providers/http/watchlist.service';
import { StateService } from '../../../providers/state/state.service';
import { environment } from '../../../../environments/environment';
import { WebsocketService } from '../../../providers/http/websocket.service';

@Component({
  selector: 'app-watchlist',
  templateUrl: './watchlist.component.html',
  host: { '(document:click)': 'onDocClick($event)' }  
})

export class WatchlistComponent implements OnInit {

  quotes = {}
  showAddWatchlist = false;
  showRenameWatchlist = false;
  showDeleteWatchlist = false;

  watchlistAdd: string = "";
  watchlistAddError: string = "";
  watchlistRename: string = "";
  watchlistRenameError: string = "";

  watchlist: Watchlist = null;
  watchlists: Watchlist[];
  activeWatchlistId: number = 0;  
  watchlistEditState = false;
  watchlistSettingsActive = false;
  sortOptions: SortablejsOptions = { animation: 150, handle: ".drag-handle" };

  //
  // Construct...
  //
  constructor(private http: HttpClient, private _eref: ElementRef, private websocketService: WebsocketService, private watchlistService: WatchlistService, private stateService: StateService) { }

  //
  // On Init...
  //
  ngOnInit() {
      
    // Set the active watchlist id
    this.activeWatchlistId = this.stateService.GetActiveWatchlistId();

    // Load up cached quotes
    this.quotes = this.stateService.GetQuotes();

    // Load watchlist from cache
    this.watchlist = this.stateService.GetActiveWatchlist();

    // Load up all watchlists
    this.getAllWatchlists();

    // Watch for changes on the watchlist order.
    this.sortOptions.onUpdate = (event: any) => {

      var ids = [];

      for(let i = 0; i < event.to.getElementsByTagName("li").length; i++)
      {
        ids.push(event.to.getElementsByTagName("li")[i].id);
      }

      // This is the new list. (watchlist_symbols)
      console.log(ids);
    };
        
    // Subscribe to data updates from the quotes - Market Quotes
    this.websocketService.quotePushData.subscribe(data => {
      this.quotes[data.symbol] = data;
    }); 
  }

  //
  // Get all watchlists from server
  //
  getAllWatchlists() 
  {
    // Get all watch lists
    this.watchlistService.get().subscribe((data) => {
      this.watchlists = data;

      if(! this.activeWatchlistId)
      {
        this.activeWatchlistId = data[0].Id;
      }

      this.setActiveWatchlist();
    });
  }

  //
  // Get active watchlist from a click
  //
  setActiveListClick(watchlist : Watchlist)
  {
    this.activeWatchlistId = watchlist.Id;
    this.setActiveWatchlist();
    this.watchlistSettingsActive = false;
  }

  //
  // Get active watchlist
  //
  setActiveWatchlist() 
  {
    // Loop through and find the watchlist by id.
    for(let i = 0; i < this.watchlists.length; i++)
    {
      if(this.watchlists[i].Id == this.activeWatchlistId)
      {
        this.watchlist = this.watchlists[i];
        this.watchlistRename = this.watchlist.Name;
        this.stateService.SetActiveWatchlist(this.watchlist);
        break; 
      }
    }
  }

  //
  // onSearchTypeAheadClick() - Add a symbol to a watch list.
  //
  onSearchTypeAheadClick(symbol: Symbol) {
    // Send request to the server.
    this.watchlistService.addSymbolByWatchlistId(this.watchlist.Id, symbol.Id).subscribe((data) => {
      this.getAllWatchlists();
    });
  }

  //
  // On watchlist settings click.
  //
  onWatchlistSettingsClick() {

    if(this.watchlistSettingsActive)
    {
      this.watchlistSettingsActive = false;
    } else
    {
      this.watchlistSettingsActive = true;      
    }
    
  } 

  //
  // Do edit watchlist
  //
  onEditWatchList() {
    this.watchlistEditState = true;
    this.watchlistSettingsActive = false;
  }

  //
  // Do static watchlist
  //
  onEditWatchListDone() {
    this.watchlistEditState = false;
    this.watchlistSettingsActive = false;
  }

  //
  // Click anywhere on the screen.
  //
  onDocClick(event) {
   
    // Remove active buttons
    if(! this._eref.nativeElement.contains(event.target))
    {
      this.watchlistSettingsActive = false;
    }
  }

  //
  // On delete watchlist.
  //
  onDeleteWatchlist(form: NgForm)
  {
    // Store this at the server.
    this.watchlistService.delete(this.watchlist.Id).subscribe((data) => {
      this.activeWatchlistId = 0;
      this.getAllWatchlists();
    });    

    this.showDeleteWatchlist = false;
    this.watchlistSettingsActive = false;    
  }

  //
  // Add a new watchlist.
  //
  onWatchlistAddSubmit(form: NgForm)
  {
    this.watchlistAddError = "";

    // Store this at the server.
    this.watchlistService.create(this.watchlistAdd).subscribe(

    // Syccess
    (data) => {
      this.watchlists.push(data);
      this.activeWatchlistId = data.Id;
      this.setActiveWatchlist();
      this.watchlistEditState = true;

      this.watchlistAdd = "";
      this.showAddWatchlist = false;
      this.watchlistSettingsActive = false;       
    },

    // Error
    (err: HttpErrorResponse) => {
      this.watchlistAddError = err.error.error;
    });  
  }

  //
  // Rename the active watchlist.
  //
  onWatchlistRenameSubmit(form: NgForm)
  {
    this.watchlistRenameError = "";

    // Store this at the server.
    this.watchlistService.update(this.watchlist.Id, this.watchlistRename).subscribe(

    // Success
    (data) => {
      this.watchlist.Name = this.watchlistRename;
      this.watchlistSettingsActive = false;
      this.showRenameWatchlist = false;
    },

    // Error
    (err: HttpErrorResponse) => {
      this.watchlistRenameError = err.error.error;
    });
  }

}

/* End File */