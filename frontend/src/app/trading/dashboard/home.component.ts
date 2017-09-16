import { Component, OnInit } from '@angular/core';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-dashboard',
  templateUrl: './home.component.html'
})
export class DashboardComponent implements OnInit {

  constructor() { }

  ngOnInit() {
    
    console.log(environment.app_server);
    
  }

}
