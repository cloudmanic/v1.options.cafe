//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, Input } from '@angular/core';
import { Order } from '../../../models/order';

@Component({
  selector: 'app-trading-orders-multileg',
  templateUrl: './orders.multileg.component.html'
})

export class OrdersMultiLegComponent {
  
  @Input() order: Order;

}

/* End File */