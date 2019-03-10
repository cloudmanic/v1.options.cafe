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
import { UpgradeComponent } from './settings/account/upgrade/upgrade.component';
import { ExpiredComponent } from './settings/account/expired/expired.component';
import { CreditCardComponent } from './settings/account/upgrade/credit-card/credit-card.component';

// Backtest
import { BacktestHomeComponent } from './backtest/home/home.component';

// Reports
import { CustomReportsComponent } from './reports/custom-reports/custom-reports.component';
import { AccountHistoryComponent } from './reports/account-history/account-history.component';
import { AccountSummaryComponent } from './reports/account-summary/account-summary.component';
import { ProfitLossComponent } from 'app/reports/custom-reports/profit-loss/profit-loss.component';
import { AccountReturnsComponent } from 'app/reports/custom-reports/account-returns/account-returns.component';
import { AccountCashComponent } from './reports/custom-reports/account-cash/account-cash.component';
import { AccountValuesComponent } from './reports/custom-reports/account-values/account-values.component';

// Trading
import { TradesComponent } from './trading/trades/home.component';
import { DashboardComponent } from './trading/dashboard/home.component';

// Screener
import { ScreenerComponent } from './trading/screener/home.component';
import { AddEditComponent as ScreenerAddEditComponent } from './trading/screener/add-edit/add-edit.component';

// Research
import { SymbolComponent } from './research/symbol/symbol.component';

// Centcom
import { CoreComponent as CentcomCoreComponent } from './centcom/layouts/core/core.component';
import { UsersComponent as CentcomUsersComponent } from './centcom/users/users.component';

// Routes
const appRoutes: Routes = [
	// Auth
	{ path: 'login', component: AuthLoginComponent },
	{ path: 'logout', redirectTo: 'login' },
	{ path: 'register', component: AuthRegisterComponent },
	{ path: 'broker-select', component: AuthBrokerSelectComponent },
	{ path: 'forgot-password', component: AuthForgotPasswordComponent },
	{ path: 'reset-password', component: AuthResetPasswordComponent },

	// Upgrade
	{ path: 'settings/account/upgrade', component: UpgradeComponent },
	{ path: 'settings/account/expired', component: ExpiredComponent },
	{ path: 'settings/account/upgrade/credit-card', component: CreditCardComponent },

	// Core App
	{
		path: '', component: LayoutCoreComponent, children: [

			// Trades
			{ path: '', component: DashboardComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'dashboard', action: '' } },
			{ path: 'trades', component: TradesComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'trades', action: '' } },
			{ path: 'screener', component: ScreenerComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'screener', action: 'list' } },
			{ path: 'screener/add', component: ScreenerAddEditComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'screener', action: 'add-edit' } },
			{ path: 'screener/edit/:id', component: ScreenerAddEditComponent, canActivate: [AuthGuard], data: { section: 'trading', subSection: 'screener', action: 'add-edit' } },

			// Reports
			{ path: 'reports/account-summary', component: AccountSummaryComponent, canActivate: [AuthGuard], data: { section: 'reports', subSection: 'account-summary', action: '' } },
			{ path: 'reports/account-history', component: AccountHistoryComponent, canActivate: [AuthGuard], data: { section: 'reports', subSection: 'account-history', action: '' } },
			{ path: 'reports/custom', component: CustomReportsComponent, canActivate: [AuthGuard], data: { section: 'reports', subSection: 'custom', action: '' } },
			{ path: 'reports/custom/profit-loss', component: ProfitLossComponent, canActivate: [AuthGuard], data: { section: 'reports', subSection: 'custom', action: '' } },
			{ path: 'reports/custom/account-cash', component: AccountCashComponent, canActivate: [AuthGuard], data: { section: 'reports', subSection: 'custom', action: '' } },
			{ path: 'reports/custom/account-values', component: AccountValuesComponent, canActivate: [AuthGuard], data: { section: 'reports', subSection: 'custom', action: '' } },
			{ path: 'reports/custom/account-returns', component: AccountReturnsComponent, canActivate: [AuthGuard], data: { section: 'reports', subSection: 'custom', action: '' } },

			// Backtest
			{ path: 'backtest', component: BacktestHomeComponent, canActivate: [AuthGuard], data: { section: 'backtest', subSection: 'dashboard', action: '' } },

			// Settings
			{ path: 'settings/account', component: AccountComponent, canActivate: [AuthGuard], data: { section: 'settings', subSection: 'account', action: '' } },
			{ path: 'settings/trading', component: TradingComponent, canActivate: [AuthGuard], data: { section: 'settings', subSection: 'trading', action: '' } },
			{ path: 'settings/brokers', component: BrokersComponent, canActivate: [AuthGuard], data: { section: 'settings', subSection: 'brokers', action: '' } },

			// Research
			{ path: 'research/symbol', component: SymbolComponent, canActivate: [AuthGuard], data: { section: 'research', subSection: 'symbol', action: '' } },


		]
	},

	// Centcom App
	{
		path: 'centcom', component: CentcomCoreComponent, children: [
			{ path: 'users', component: CentcomUsersComponent, canActivate: [AuthGuard], data: {} },
		]
	},


	// Otherwise redirect to home
	{ path: '**', redirectTo: '' }
];


export const Routing = RouterModule.forRoot(appRoutes);
