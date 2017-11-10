//
// Date: 11/10/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Symbol } from '../../models/symbol';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-typeahead-symbols',
  templateUrl: './typeahead-symbols.component.html'  
})

export class TypeaheadSymbolsComponent implements OnInit {

  @Output('selected') symbolSelected = new EventEmitter<Symbol>();

  typeAheadList: Symbol[];
  typeAheadShow = false;  

  //
  // Construct.
  //
  constructor(private http: HttpClient) { }

  //
  // On Init.
  //
  ngOnInit() {}

  //
  // On item selected from the type ahead.
  //
  onSelected(symbol: Symbol) {
    this.symbolSelected.emit(symbol)
  }

  //
  // On search...
  //
  onSearchKeyUp(event) {

    // Send search to backend.
    if(event.target.value.length > 0)
    {
      this.typeAheadShow = true;
    } else
    {
      this.typeAheadList = []
      this.typeAheadShow = false;
      return false; // No ajax call needed.
    }

    // Send this search even to the server to get results.
    this.http.get<Symbol[]>(environment.app_server + '/api/v1/symbols?search=' + event.target.value).subscribe(
      
      // Success
      (symbols: Symbol[]) => {
        if(symbols)
        {
          this.typeAheadList = symbols;
          this.typeAheadShow = true;
        } else
        {
          this.typeAheadList = []
          this.typeAheadShow = false;
        }
      },
      
      // Error
      (err: HttpErrorResponse) => {

        if(err.error instanceof Error) 
        {
          // A client-side or network error occurred. Handle it accordingly.
          console.log('An error occurred:', err.error.message);
        } else 
        { 
          // Print error message
          var json = JSON.parse(err.error); // Bug....Angular 4.4.4
          console.log(json)
        }
      }
    );
  }

}

/* End File */