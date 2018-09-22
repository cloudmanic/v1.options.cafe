//
// Date: 9/19/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { Broker } from '../../models/broker';
import { BrokerAccount } from '../../models/broker-account';
import { BrokerService } from '../../providers/http/broker.service';

@Component({
  selector: 'app-brokers',
  templateUrl: './brokers.component.html',
  styleUrls: []
})

export class BrokersComponent implements OnInit 
{
  brokers: Broker[] = []

  //
  // Construct.
  //
  constructor(private brokerService: BrokerService) 
  { 
    // Load data.
    this.getBrokers();
  }

  //
  // NgInit
  //
  ngOnInit() {}

  //
  // Get brokers.
  //
  getBrokers()
  {
    // Ajax call to get brokers.
    this.brokerService.get().subscribe((res) => {
      this.brokers = res;

      for (let i = 0; i < this.brokers.length; i++)
      {
        if (this.brokers[i].BrokerAccounts.length > 0) 
        {
          this.brokers[i].SettingsActiveBrokerAccount = this.brokers[i].BrokerAccounts[0];
        }
      }
    });
  }

  //
  // Broker account click
  //
  brokerAccountClick(broker: Broker, row: BrokerAccount)
  {
    broker.SettingsActiveBrokerAccount = row;
  }

  //
  // Return a CSS for the logo of this broker.
  //
  getLogoClass(row: Broker) : string 
  {
    let cssClass: string = '';

    switch(row.Name)
    {
      case 'Tradier':
        cssClass = 'logo-tradier';
      break;

      case 'Tradier Sandbox':
        cssClass = 'logo-tradier';
      break;
    }

    return cssClass;
  }
}

/* End File */