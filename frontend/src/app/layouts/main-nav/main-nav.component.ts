declare var groove: any;

import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-main-nav',
  templateUrl: './main-nav.component.html'
})

export class MainNavComponent implements OnInit {

  constructor() { }

  ngOnInit() {
  }

  //
  // Clicked on help.
  //
  onHelpClick() { 
    groove.widget('open');
  }

}
