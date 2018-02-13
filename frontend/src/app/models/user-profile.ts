//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

export class UserProfile 
{
  
  //
  // Construct...
  //
  constructor(
    public Id: string,
    public Name: string
  ){}
  
  //
  // Build user profile for emitting to the app.
  //
  public static buildForEmit(data) : UserProfile  {
    
    let user = new UserProfile(data.Id, data.Name); 
            
    return user;
  }
  
}

/* End File */