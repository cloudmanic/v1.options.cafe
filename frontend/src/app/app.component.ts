import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html'
})

export class AppComponent 
{
  title = 'app';

  //
  // Constructor
  //
  constructor(private router: Router) 
  {
    // redirect
    let redt = localStorage.getItem('redirect');

    if(redt) 
    {
      localStorage.removeItem('redirect');
      this.router.navigate([redt]);
    }    
  }
}

/* End File */