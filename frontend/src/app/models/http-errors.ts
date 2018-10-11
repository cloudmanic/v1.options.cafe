//
// Date: 10/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

export class HttpErrors 
{
  FirstName: string;
  LastName: string;
  Email: string;
  Phone: string;

  //
  // Build from JSON.
  //
  fromJson(json: Object): HttpErrors 
  {
    let result = new HttpErrors();

    if (!json) {
      return result;
    }

    // Set data.
    if(typeof result.FirstName == "undefined")
    {
      result.FirstName = json["first_name"];
    }

    if (typeof result.LastName == "undefined") 
    {
      result.LastName = json["last_name"];
    }

    if (typeof result.Email == "undefined") 
    {
      result.Email = json["email"];
    }

    if (typeof result.Phone == "undefined") 
    {
      result.Phone = json["phone"];
    }

    // Return happy
    return result;
  }
}

/* End File */