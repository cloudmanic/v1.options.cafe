//
// Date: 10/11/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';

export class Invoice 
{
  Amount: number;
  Date: Date;
  InvoiceUrl: string;
  PaymentMethod: string;
  Transaction: string; 

  //
  // Build from json list.
  //
  fromJsonList(json: Object[]): Invoice[]
  {
    let invoices: Invoice[] = [];

    if (!json) 
    {
      return invoices;
    }

    for (let i = 0; i < json.length; i++)
    {
      invoices.push(this.fromJson(json[i]));      
    }

    return invoices;
  }

  //
  // Build from JSON.
  //
  fromJson(json: Object): Invoice 
  {
    let result = new Invoice();

    if (!json) {
      return result;
    }

    // Set data.
    result.Amount = json["amount"];
    result.Date = moment(json["date"]).toDate();
    result.InvoiceUrl = json["invoice_url"];
    result.PaymentMethod = json["payment_method"];
    result.Transaction = json["transaction"];   

    // Return happy
    return result;
  }
}



/* End File */
