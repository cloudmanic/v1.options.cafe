
//
// Date: 11/11/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

import { HttpRequest, HttpHandler, HttpEvent, HttpInterceptor } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';

export class TokenInterceptor implements HttpInterceptor 
{
  //
  // Construct
  //
  constructor() {}
  
  //
  // Run intercept
  //
  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    
    let access_token = localStorage.getItem('access_token');

    // Add in the headers for the API.
    request = request.clone({ setHeaders: { Authorization: `Bearer ${access_token}` } });

    // Let the request continue.
    return next.handle(request);
  }
}

/* End File */