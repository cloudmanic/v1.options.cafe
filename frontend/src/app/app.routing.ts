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
import { SettingsComponent } from './settings/settings.component';

// Backtest
import { BacktestHomeComponent } from './backtest/home/home.component';

// Reports
import { ReportsHomeComponent } from './reports/home/home.component';

// Trading
import { TradesComponent } from './trading/trades/home.component';
import { ScreenerComponent } from './trading/screener/home.component';
import { DashboardComponent } from './trading/dashboard/home.component';

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
    
    { path: '', component: DashboardComponent, canActivate: [ AuthGuard ] },
    { path: 'screener', component: ScreenerComponent, canActivate: [ AuthGuard ] },
    { path: 'trades', component: TradesComponent, canActivate: [ AuthGuard ] },
    { path: 'reports', component: ReportsHomeComponent, canActivate: [ AuthGuard ] },
    { path: 'backtest', component: BacktestHomeComponent, canActivate: [ AuthGuard ] },
    { path: 'settings', component: SettingsComponent, canActivate: [AuthGuard] }
    
  ] },
  
  // Otherwise redirect to home
  { path: '**', redirectTo: '' }  
];


export const Routing = RouterModule.forRoot(appRoutes);