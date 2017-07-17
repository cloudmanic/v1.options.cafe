import { Component, OnInit } from '@angular/core';
import { AppService } from '../../providers/websocket/app.service';

@Component({
  selector: 'app-auth-register',
  templateUrl: './home.component.html'
})
export class AuthRegisterComponent implements OnInit {

  constructor(private appService: AppService) { }

  ngOnInit() {
  }

}
