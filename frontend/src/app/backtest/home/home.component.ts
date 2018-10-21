import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-backtest-home',
  templateUrl: './home.component.html'
})

export class BacktestHomeComponent implements OnInit 
{

  constructor() { }

  ngOnInit() {}


  //
  // Share on twitter
  //
  twitterShare()
  {
    let tweet = "Options backtesting is coming soon to Options Cafe!";
    window.open('https://twitter.com/share?text=' + tweet + '&via=options_cafe&url=https://options.cafe&hashtags=OptionsTrading', '', 'menubar=no, toolbar = no, resizable = yes, scrollbars = yes, height = 600, width = 600');
  }

  //
  // Share on facebook
  //
  facebookShare() 
  {
    window.open('https://www.facebook.com/sharer/sharer.php?u=https://options.cafe', '', 'menubar=no, toolbar = no, resizable = yes, scrollbars = yes, height = 600, width = 600');
  }

  //
  // Share on google
  //
  googleShare() {
    window.open('https://plus.google.com/share?url=https://options.cafe', '', 'menubar=no, toolbar = no, resizable = yes, scrollbars = yes, height = 600, width = 600');
  }    
}
