//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, Input } from '@angular/core';
import { Order } from '../../../models/order';

@Component({
  selector: 'app-trading-orders-option',
  templateUrl: './orders.option.component.html'
})

export class OrdersOptionComponent {
  
  @Input() order: Order;
  @Input() quotes: {};

}

/* End File */