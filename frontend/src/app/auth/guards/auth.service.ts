//
// Date: 9/16/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { Injectable } from '@angular/core';
import { Router, CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';

@Injectable()
export class AuthGuard implements CanActivate {

  //
  // Construct.
  //
  constructor(private router: Router) {}

  //
  // Is the user allowed to use this page?
  //
  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot) {
    
    // Make sure we have an access token.
    if(localStorage.getItem('access_token') && localStorage.getItem('user_id')) {
      return true;
    }

    // Not logged in so redirect to login page with the return url
    this.router.navigate(['/login'], { queryParams: { returnUrl: state.url }});
    return false;
        
  }
}