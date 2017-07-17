import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { AppComponent } from './app.component';

// Providers
import { AppService } from './providers/websocket/app.service';

// Layout
import { SidebarComponent } from './layouts/sidebar/sidebar.component';
import { MainNavComponent } from './layouts/main-nav/main-nav.component';

// Auth
import { AuthLayoutComponent } from './auth/layout/home.component';
import { AuthLoginComponent } from './auth/login/home.component';
import { AuthRegisterComponent } from './auth/register/home.component';
import { AuthBrokerSelectComponent } from './auth/broker-select/home.component';
import { AuthResetPasswordComponent } from './auth/reset-password/home.component';

// Backtest
import { BacktestLayoutComponent } from './backtest/layout/layout.component';
import { BacktestHomeComponent } from './backtest/home/home.component';

// Reports
import { ReportsLayoutComponent } from './reports/layout/layout.component';
import { ReportsHomeComponent } from './reports/home/home.component';

// Trading
import { TradingLayoutComponent } from './trading/layout/layout.component';
import { TradesComponent } from './trading/trades/home.component';
import { ScreenerComponent } from './trading/screener/home.component';
import { DashboardComponent } from './trading/dashboard/home.component';


// Routes
const appRoutes: Routes = [
  // Auth
  { path: 'login', component: AuthLoginComponent },  
  { path: 'register', component: AuthRegisterComponent },
  { path: 'broker_select', component: AuthBrokerSelectComponent }, 
  { path: 'reset_password', component: AuthResetPasswordComponent },   
  
  // Trading
  { path: '', component: DashboardComponent },
  { path: 'dashboard', component: DashboardComponent },  
  { path: 'screener', component: ScreenerComponent },
  { path: 'trades', component: TradesComponent },
    
  // Reports
  { path: 'reports', component: ReportsHomeComponent },
  
  // Backtest
  { path: 'backtest', component: BacktestHomeComponent }
];

@NgModule({
  declarations: [
    AppComponent,
    
    // Layout
    SidebarComponent,
    MainNavComponent,    
    
    // Auth
    AuthLayoutComponent,
    AuthLoginComponent,
    AuthRegisterComponent,
    AuthBrokerSelectComponent,
    AuthResetPasswordComponent,
    
    // Backtest
    BacktestLayoutComponent,
    BacktestHomeComponent,

    // Reports
    ReportsLayoutComponent,
    ReportsHomeComponent,
    
    // Trading
    TradingLayoutComponent,
    TradesComponent,
    ScreenerComponent,
    DashboardComponent
  ],
  imports: [
    BrowserModule,
    RouterModule.forRoot(appRoutes)
  ],
  providers: [ AppService ],
  bootstrap: [AppComponent]
})
export class AppModule { }
