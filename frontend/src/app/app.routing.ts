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
  { path: 'logout', redirectTo: 'login' },    
  { path: 'register', component: AuthRegisterComponent },
  { path: 'broker-select', component: AuthBrokerSelectComponent }, 
  { path: 'forgot-password', component: AuthForgotPasswordComponent },   
  { path: 'reset-password', component: AuthResetPasswordComponent }, 
  
  // Trading
  { path: 'dashboard', component: DashboardComponent, canActivate: [ AuthGuard ] },  
  { path: 'screener', component: ScreenerComponent, canActivate: [ AuthGuard ] },
  { path: 'trades', component: TradesComponent, canActivate: [ AuthGuard ] },
    
  // Reports
  { path: 'reports', component: ReportsHomeComponent, canActivate: [ AuthGuard ] },
  
  // Backtest
  { path: 'backtest', component: BacktestHomeComponent, canActivate: [ AuthGuard ] },
  
  // Otherwise redirect to home
  { path: '**', redirectTo: 'dashboard' }  
];


export const Routing = RouterModule.forRoot(appRoutes);