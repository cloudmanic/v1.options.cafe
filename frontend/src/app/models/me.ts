//
// Date: 10/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

export class Me 
{
  Id: number;
  FirstName: string;
  LastName: string;
  Email: string;
  Phone: string;
  Address: string;
  City: string;
  State: string;
  Zip: string;
  Country: string;
  GoogleSubId: string;

  //
  // Build from JSON.
  //
  fromJson(json: Object) : Me 
  {
    let result = new Me();

    if (!json) {
      return result;
    }

    // Set data.
    result.Id = json["id"];
    result.FirstName = json["first_name"];
    result.LastName = json["last_name"];
    result.Email = json["email"];
    result.Phone = json["phone"];
    result.Address = json["address"];
    result.City = json["city"];
    result.State = json["state"];
    result.Zip = json["zip"];
    result.Country = json["country"];
    result.GoogleSubId = json["google_sub_id"];

    // Return happy
    return result;
  }

  //
  // Set from object
  //
  setFromObj(obj: Me) : Me 
  {
    let result = new Me();

    // Set data.
    result.Id = obj.Id;
    result.FirstName = obj.FirstName;
    result.LastName = obj.LastName;
    result.Email = obj.Email;
    result.Phone = obj.Phone;
    result.Address = obj.Address;
    result.City = obj.City;
    result.State = obj.State;
    result.Zip = obj.Zip;
    result.Country = obj.Country;
    result.GoogleSubId = obj.GoogleSubId;

    // Return happy
    return result;    
  }
}

/* End File */
