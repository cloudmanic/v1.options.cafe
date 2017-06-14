import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { AppComponent } from './app.component';

// Layout
import { SidebarComponent } from './layouts/sidebar/sidebar.component';
import { MainNavComponent } from './layouts/main-nav/main-nav.component';

import { BacktestingComponent } from './backtesting/backtesting.component';
import { ReportsComponent } from './reports/reports.component';

// Trading
import { TradingLayoutComponent } from './trading/layout/layout.component';
import { TradingDashboardComponent } from './trading/dashboard/dashboard.component';

const appRoutes: Routes = [
  { path: '', component: TradingDashboardComponent },
  { path: 'reports', component: ReportsComponent },
  { path: 'backtesting', component: BacktestingComponent }
];

@NgModule({
  declarations: [
    AppComponent,
    
    // Layout
    SidebarComponent,
    MainNavComponent,    
    
    BacktestingComponent,
    ReportsComponent,
    
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
