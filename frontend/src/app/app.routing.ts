import { Routes, RouterModule } from '@angular/router';

// Guards
import { AuthGuard } from './auth/guards/auth.service';

// Auth
import { AuthLayoutComponent } from './auth/layout/home.component';
import { AuthLoginComponent } from './auth/login/home.component';
import { AuthRegisterComponent } from './auth/register/home.component';
import { AuthBrokerSelectComponent } from './auth/broker-select/home.component';
import { AuthResetPasswordComponent } from './auth/reset-password/home.component';
import { AuthForgotPasswordComponent } from './auth/forgot-password/home.component';

// Layouts
import { LayoutCoreComponent } from './layouts/core/core.component';

// Setting
import { AccountComponent } from './settings/account/account.component';
import { TradingComponent } from './settings/trading/trading.component';
import { BrokersComponent } from './settings/brokers/brokers.component';

// Backtest
import { BacktestHomeComponent } from './backtest/home/home.component';

// Reports
import { AccountHistoryComponent } from './reports/account-history/account-history.component';
import { AccountSummaryComponent } from './reports/account-summary/account-summary.component';

// Trading
import { TradesComponent } from './trading/trades/home.component';
import { DashboardComponent } from './trading/dashboard/home.component';

// Screener
import { ScreenerComponent } from './trading/screener/home.component';
import { AddEditComponent as ScreenerAddEditComponent } from './trading/screener/add-edit/add-edit.component';

// Routes
const appRoutes: Routes = [
  // Auth
  { path: 'login', component: AuthLoginComponent }, 
  { path: 'logout', redirectTo: 'login' },    
  { path: 'register', component: AuthRegisterComponent },
  { path: 'broker-select', component: AuthBrokerSelectComponent }, 
  { path: 'forgot-password', component: AuthForgotPasswordComponent },   
  { path: 'reset-password', component: AuthResetPasswordComponent }, 
  
  // Core App
  { path: '', component: LayoutCoreComponent, children: [
    
    // Trades
    { path: '', component: DashboardComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'dashboard', action: '' } },
    { path: 'trades', component: TradesComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'trades', action: '' } },    
    { path: 'screener', component: ScreenerComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'screener', action: 'list' } },
    { path: 'screener/add', component: ScreenerAddEditComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'screener', action: 'add-edit' } },
    { path: 'screener/edit/:id', component: ScreenerAddEditComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'screener', action: 'add-edit' } },

    // Reports
    { path: 'reports/account-summary', component: AccountSummaryComponent, canActivate: [AuthGuard], data: { section: 'reports', subSection: 'account-summary', action: '' } },
    { path: 'reports/account-history', component: AccountHistoryComponent, canActivate: [AuthGuard], data: { section: 'reports', subSection: 'account-history', action: '' } },

    // Backtest 
    { path: 'backtest', component: BacktestHomeComponent, canActivate: [AuthGuard], data: { section: 'backtest', subSection: 'dashboard', action: '' } },

    // Settings
    { path: 'settings/account', component: AccountComponent, canActivate: [AuthGuard], data: { section: 'settings', subSection: 'account', action: '' } },
    { path: 'settings/trading', component: TradingComponent, canActivate: [AuthGuard], data: { section: 'settings', subSection: 'trading', action: '' } },
    { path: 'settings/brokers', component: BrokersComponent, canActivate: [AuthGuard], data: { section: 'settings', subSection: 'brokers', action: '' } }
    
  ] },
  
  // Otherwise redirect to home
  { path: '**', redirectTo: '' }  
];


export const Routing = RouterModule.forRoot(appRoutes);