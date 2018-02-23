import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule }   from '@angular/forms';
import { Routes, RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';
import { Routing } from './app.routing';
import { SortablejsModule } from 'angular-sortablejs';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { TokenInterceptor } from './providers/http/token.interceptor';

// Shared
import { PagingComponent } from './shared/paging/paging.component';
import { DialogComponent } from './shared/dialog/dialog.component';

// Providers
import { AuthGuard } from './auth/guards/auth.service';
import { AppService } from './providers/websocket/app.service';
import { QuoteService } from './providers/websocket/quote.service';
import { StateService } from './providers/state/state.service';

// Providers - http
import { BrokerService } from './providers/http/broker.service';
import { SymbolService } from './providers/http/symbol.service';
import { WatchlistService } from './providers/http/watchlist.service';
import { TradeGroupService } from './providers/http/trade-group.service';

// Layout
import { SidebarComponent } from './layouts/sidebar/sidebar.component';
import { MainNavComponent } from './layouts/main-nav/main-nav.component';
import { LayoutCoreComponent } from './layouts/core/core.component';

// Auth
import { AuthLayoutComponent } from './auth/layout/home.component';
import { AuthLoginComponent } from './auth/login/home.component';
import { AuthRegisterComponent } from './auth/register/home.component';
import { AuthBrokerSelectComponent } from './auth/broker-select/home.component';
import { AuthResetPasswordComponent } from './auth/reset-password/home.component';
import { AuthForgotPasswordComponent } from './auth/forgot-password/home.component';

// Backtest
import { BacktestSubnavComponent } from './backtest/sub-nav/subnav.component';
import { BacktestHomeComponent } from './backtest/home/home.component';

// Reports
import { ReportsSubnavComponent } from './reports/sub-nav/subnav.component';
import { ReportsHomeComponent } from './reports/home/home.component';

// Trading
import { IvrComponent } from './trading/dashboard/ivr/ivr.component';
import { SubnavComponent } from './trading/sub-nav/subnav.component';
import { TradesComponent } from './trading/trades/home.component';
import { WatchlistComponent } from './trading/dashboard/watchlist/watchlist.component';
import { ScreenerComponent } from './trading/screener/home.component';
import { DashboardComponent } from './trading/dashboard/home.component';
import { OrdersComponent } from './trading/dashboard/orders/orders.component';
import { PositionsComponent } from './trading/dashboard/positions/positions.component';
import { PositionComponent } from './trading/dashboard/positions/position.component';
import { MarketQuotesComponent } from './trading/dashboard/market-quotes/market-quotes.component';
import { DashboardChartComponent } from './trading/dashboard/dashboard-chart/dashboard-chart.component';
import { TypeaheadSymbolsComponent } from './shared/typeahead-symbols/typeahead-symbols.component';

@NgModule({
  declarations: [
    
    AppComponent,
    
    // Layout
    SidebarComponent,
    MainNavComponent,
    LayoutCoreComponent,    
    
    // Shared
    PagingComponent,
    DialogComponent,    

    // Auth
    AuthLayoutComponent,
    AuthLoginComponent,
    AuthRegisterComponent,
    AuthBrokerSelectComponent,
    AuthResetPasswordComponent,
    AuthForgotPasswordComponent,
    
    // Backtest
    BacktestSubnavComponent,
    BacktestHomeComponent,

    // Reports
    ReportsSubnavComponent,
    ReportsHomeComponent,
    
    // Trading
    IvrComponent,
    OrdersComponent,
    SubnavComponent,
    TradesComponent,
    ScreenerComponent,
    DashboardComponent,
    WatchlistComponent,
    PositionsComponent,
    PositionComponent,
    MarketQuotesComponent,
    DashboardChartComponent,
    TypeaheadSymbolsComponent,
  ],
  
  imports: [
    Routing,
    FormsModule,
    BrowserModule,
    SortablejsModule,
    HttpClientModule   
  ],
  
  providers: [ 
    AppService, 
    QuoteService, 
    AuthGuard, 
    SymbolService,
    BrokerService,
    TradeGroupService,
    StateService,
    WatchlistService,
    { provide: HTTP_INTERCEPTORS, useClass: TokenInterceptor, multi: true }    
  ],
  
  bootstrap: [AppComponent]
})

export class AppModule { }