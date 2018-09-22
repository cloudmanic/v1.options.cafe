//
// Date: 9/19/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { Broker } from '../../models/broker';
import { BrokerAccount } from '../../models/broker-account';
import { BrokerService } from '../../providers/http/broker.service';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-brokers',
  templateUrl: './brokers.component.html',
  styleUrls: []
})

export class BrokersComponent implements OnInit 
{
  brokers: Broker[] = []
  showAddEditBroker: boolean = false;

  // Add Broker Stuff
  addBrokerType: string = "Tradier";
  addBrokerError: string = "";
  addBrokerDisplayName: string = "";

  //
  // Construct.
  //
  constructor(private brokerService: BrokerService) 
  { 
    // TODO: Do some sort of notice with this.`
    localStorage.removeItem('broker_new_id');

    // Load data.
    this.getBrokers();
  }

  //
  // NgInit
  //
  ngOnInit() {}

  //
  // Add broker
  //
  addBroker()
  {
    if(this.addBrokerDisplayName.length <= 0)
    {
      this.addBrokerError = "A broker display name is required.";
      return;
    } else
    {
      this.addBrokerError = "";
    }

    // Ajax call to add broker.
    this.brokerService.create(this.addBrokerType, this.addBrokerDisplayName).subscribe((res) => {

      // Set redirect for after auth with brpker
      localStorage.setItem('redirect', '/settings/brokers');
      localStorage.setItem('broker_new_id', String(res.Id));

      // Switch based on broker selected - Redirect to login to broker and get access token.
      switch (this.addBrokerType) 
      {
        case 'Tradier':
          window.location.href = environment.app_server + '/tradier/authorize?user=' + localStorage.getItem('user_id') + '&broker_id=' + res.Id;
        break;
      }

    });    
  }

  //
  // Show add broker
  //
  showAddBrokerPopup()
  {
    this.addBrokerError = "";
    this.addBrokerDisplayName = "";
    this.showAddEditBroker = true;
  }

  //
  // Close add broker
  //
  closeShowAddEditBroker() 
  {
    this.showAddEditBroker = false;    
  }

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