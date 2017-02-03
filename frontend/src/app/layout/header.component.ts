import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { BrokerService } from '../services/broker.service';
import { QuoteService } from '../services/quote.service';
import { UserProfile } from '../contracts/user-profile';

@Component({
  selector: 'oc-header',
  templateUrl: './header.component.html'
})

export class HeaderComponent implements OnInit {
  quotes = {};
  userProfile: UserProfile;

  //
  // Constructor....
  //
  constructor(private broker: BrokerService, private quotesService: QuoteService, private changeDetect: ChangeDetectorRef) { }

  //
  // OnInit....
  //
  ngOnInit() {

    // Subscribe to data updates from the broker - Market Status
    this.broker.userProfilePushData.subscribe(data => {
      this.userProfile = data;
      this.changeDetect.detectChanges();
    });
    
    // Subscribe to data updates from the quotes - Market Quotes
    this.quotesService.marketQuotePushData.subscribe(data => {
      
      this.quotes[data.symbol] = data;
      this.changeDetect.detectChanges();
      
    });
        
  }

}

/* End File */