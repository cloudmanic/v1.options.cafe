import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { AppComponent } from './app.component';

// Layout
import { SidebarComponent } from './layouts/sidebar/sidebar.component';
import { MainNavComponent } from './layouts/main-nav/main-nav.component';

// Backtest
import { BacktestLayoutComponent } from './backtest/layout/layout.component';
import { BacktestHomeComponent } from './backtest/home/home.component';

// Reports
import { ReportsLayoutComponent } from './reports/layout/layout.component';
import { ReportsHomeComponent } from './reports/home/home.component';

// Trading
import { TradingLayoutComponent } from './trading/layout/layout.component';
import { TradingDashboardComponent } from './trading/dashboard/dashboard.component';

// Routes
const appRoutes: Routes = [
  { path: '', component: TradingDashboardComponent },
  { path: 'reports', component: ReportsHomeComponent },
  { path: 'backtest', component: BacktestHomeComponent }
];

@NgModule({
  declarations: [
    AppComponent,
    
    // Layout
    SidebarComponent,
    MainNavComponent,    
    
    // Backtest
    BacktestLayoutComponent,
    BacktestHomeComponent,

    // Reports
    ReportsLayoutComponent,
    ReportsHomeComponent,
    
    // Trading
    TradingLayoutComponent,
    TradingDashboardComponent
  ],
  imports: [
    BrowserModule,
    RouterModule.forRoot(appRoutes)
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
