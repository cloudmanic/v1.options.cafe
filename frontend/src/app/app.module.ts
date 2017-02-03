import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component';
import { QuoteService } from './services/quote.service';
import { BrokerService } from './services/broker.service';
import { HeaderComponent } from './layout/header.component';
import { FooterComponent } from './layout/footer.component';
import { AccountsComponent } from './accounts/accounts.component';
import { WatchlistComponent } from './watchlist/watchlist.component';
import { OrdersComponent } from './orders/orders.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { SidebarComponent } from './sidebar/sidebar.component';

@NgModule({
  declarations: [
    AppComponent,
    HeaderComponent,
    FooterComponent,
    AccountsComponent,
    WatchlistComponent,
    OrdersComponent,
    DashboardComponent,
    SidebarComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule
  ],
  providers: [ BrokerService, QuoteService ],
  bootstrap: [ AppComponent ]
})

export class AppModule {}
