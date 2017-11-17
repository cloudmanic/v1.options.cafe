//
// Date: 11/10/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/Rx'
import { Symbol } from '../../models/symbol';
import { SymbolService } from '../../providers/http/symbol.service';
import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'app-typeahead-symbols',
  templateUrl: './typeahead-symbols.component.html',
  host: { '(window:keydown)': 'onKeyDown($event)' }     
})

export class TypeaheadSymbolsComponent implements OnInit {

  @Output('selected') symbolSelected = new EventEmitter<Symbol>();

  searchTerm: string = ""
  activeItem: number = -1;
  typeAheadList: Symbol[];
  typeAheadShow = false;  

  //
  // Construct.
  //
  constructor(private symbolService: SymbolService) { }

  //
  // On Init.
  //
  ngOnInit() {}

  //
  // On item selected from the type ahead.
  //
  onSelected(symbol: Symbol) {
    this.symbolSelected.emit(symbol);
    this.typeAheadList = []
    this.typeAheadShow = false;
    this.activeItem = -1;
    this.searchTerm = '';        
  }

  //
  // Catch keys we we can move the items of the type ahead. 
  //
  onKeyDown(event) {

    // Nothing to do if type ahead is not showing
    if(! this.typeAheadShow)
    {
      return;
    }

    // Key down
    if(event.which === 40 || event.keyCode === 40) 
    {
      if((this.activeItem + 1) < this.typeAheadList.length)
      {
        this.activeItem++;
      }
    }

    // Key up
    if(event.which === 38 || event.keyCode === 38) 
    {
      if(this.activeItem > 0)
      {
        this.activeItem--;
      }
    }

    // Key Enter
    if(event.which === 10 || event.which === 13 || event.keyCode === 10 || event.keyCode === 13)
    {
      if(this.activeItem >= 0)
      {
        this.onSelected(this.typeAheadList[this.activeItem]);
      }
    }
  }

  //
  // On search...
  //
  onSearchKeyUp(event) {

    // Key Enter (do nothing) This mean we selected an item.
    if(event.which === 10 || event.which === 13 || event.keyCode === 10 || event.keyCode === 13)
    {
      return false;
    }

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
    this.symbolService.searchSymbols(event.target.value).subscribe((data) => {
      
      if(data)
      {
        this.typeAheadList = data;
        this.typeAheadShow = true;
      } else
      {
        this.typeAheadList = []
        this.typeAheadShow = false;
      }

    });
  }
}

/* End File */