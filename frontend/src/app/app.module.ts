import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Routes, RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';
import { Routing } from './app.routing';
import { SortablejsModule } from 'angular-sortablejs';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { TokenInterceptor } from './providers/http/token.interceptor';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { HighchartsChartComponent } from './shared/highcharts/highcharts-chart.component';
import { ServiceWorkerModule } from '@angular/service-worker';

// Shared
import { AnalyzeComponent } from './shared/analyze/analyze.component';
import { ToolTipComponent } from './shared/tool-tip/tool-tip.component';
import { PagingComponent } from './shared/paging/paging.component';
import { DialogComponent } from './shared/dialog/dialog.component';

// Providers
import { AuthGuard } from './auth/guards/auth.service';
import { StateService } from './providers/state/state.service';

// Providers - http
import { MeService } from './providers/http/me.service';
import { OptionsChainService } from './providers/http/options-chain.service';
import { TradeService } from './providers/http/trade.service';
import { ScreenerService } from './providers/http/screener.service';
import { QuotesService } from './providers/http/quotes.service';
import { BrokerService } from './providers/http/broker.service';
import { SymbolService } from './providers/http/symbol.service';
import { StatusService } from './providers/http/status.service';
import { ReportsService } from './providers/http/reports.service';
import { SettingsService } from './providers/http/settings.service';
import { WatchlistService } from './providers/http/watchlist.service';
import { TradeGroupService } from './providers/http/trade-group.service';
import { WebsocketService } from './providers/http/websocket.service';
import { AnalyzeService } from './providers/http/analyze.service';
import { BrokerEventsService } from './providers/http/broker-events.service';
import { NotificationsService } from './providers/http/notifications.service';
import { BacktestService } from './providers/http/backtest.service';

// Layout
import { SubnavComponent } from './layouts/sub-nav/subnav.component';
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
import { BacktestViewComponent } from './backtest/view/view.component';
import { BacktestSubnavComponent } from './backtest/sub-nav/subnav.component';
import { BacktestHomeComponent } from './backtest/home/home.component';
import { BacktestCreateComponent } from './backtest/create/create.component';
import { BacktestComingSoonComponent } from './backtest/coming-soon/coming-soon.component';

// Reports
import { BaseComponent } from './reports/custom-reports/base/base.component';
import { CustomReportsComponent } from './reports/custom-reports/custom-reports.component';
import { AccountHistoryComponent } from './reports/account-history/account-history.component';
import { AccountSummaryComponent } from './reports/account-summary/account-summary.component';
import { AccountReturnsComponent } from './reports/custom-reports/account-returns/account-returns.component';
import { ProfitLossComponent } from './reports/custom-reports/profit-loss/profit-loss.component';
import { AccountCashComponent } from './reports/custom-reports/account-cash/account-cash.component';
import { AccountValuesComponent } from './reports/custom-reports/account-values/account-values.component';

// Settings
import { BrokersComponent } from './settings/brokers/brokers.component';
import { TradingComponent } from './settings/trading/trading.component';
import { AccountComponent } from './settings/account/account.component';
import { PersonalInfoComponent } from './settings/account/personal-info/personal-info.component';
import { BillingHistoryComponent } from './settings/account/billing-history/billing-history.component';
import { AccountDetailsComponent } from './settings/account/account-details/account-details.component';
import { SocialComponent } from './settings/account/social/social.component';
import { UpgradeComponent } from './settings/account/upgrade/upgrade.component';
import { ExpiredComponent } from './settings/account/expired/expired.component';
import { CreditCardComponent } from './settings/account/upgrade/credit-card/credit-card.component';
import { CardComponent } from './settings/account/card/card.component';

// Trading
import { IvrComponent } from './trading/dashboard/ivr/ivr.component';
import { TradesComponent } from './trading/trades/home.component';
import { WatchlistComponent } from './trading/dashboard/watchlist/watchlist.component';
import { DashboardComponent } from './trading/dashboard/home.component';
import { OrdersComponent } from './trading/dashboard/orders/orders.component';
import { PositionsComponent } from './trading/dashboard/positions/positions.component';
import { MarketQuotesComponent } from './trading/dashboard/market-quotes/market-quotes.component';
import { DashboardChartComponent } from './trading/dashboard/dashboard-chart/dashboard-chart.component';
import { TypeaheadSymbolsComponent } from './shared/typeahead-symbols/typeahead-symbols.component';
import { TradeComponent } from './trade/trade.component';
import { TradeMultiLegComponent } from './trade/trade-multi-leg/trade-multi-leg.component';
import { DropdownSelectComponent } from './shared/dropdown-select/dropdown-select.component';
import { EquityComponent } from './trading/dashboard/positions/types/equity/equity.component';
import { TradeOptionComponent } from './trade/trade-option/trade-option.component';
import { TradeEquityComponent } from './trade/trade-equity/trade-equity.component';
import { SpreadComponent } from './trading/dashboard/positions/types/spread/spread.component';
import { OtherComponent } from './trading/dashboard/positions/types/other/other.component';
import { OptionComponent } from './trading/dashboard/positions/types/option/option.component';
import { IronCondorComponent } from './trading/dashboard/positions/types/iron-condor/iron-condor.component';
import { ReverseIronCondorComponent } from './trading/dashboard/positions/types/reverse-iron-condor/reverse-iron-condor.component';
import { LongCallButterflyComponent } from './trading/dashboard/positions/types/long-call-butterfly/long-call-butterfly.component';
import { LongPutButterflyComponent } from './trading/dashboard/positions/types/long-put-butterfly/long-put-butterfly.component';

// Screener
import { ScreenerComponent } from './trading/screener/home.component';
import { AddEditComponent as ScreenerAddEditComponent } from './trading/screener/add-edit/add-edit.component';

// Pipes
import { TableSortPipe } from './pipes/table-sort.pipe';


// Research
import { SymbolComponent } from './research/symbol/symbol.component';

// Centcom
import { CoreComponent as CentcomCoreComponent } from './centcom/layouts/core/core.component';
import { UsersComponent as CentcomUsersComponent } from './centcom/users/users.component';
import { environment } from 'environments/environment';

@NgModule({
	declarations: [
		AppComponent,
		HighchartsChartComponent,

		// Pipes
		TableSortPipe,

		// Layout
		SidebarComponent,
		MainNavComponent,
		LayoutCoreComponent,

		// Shared
		ToolTipComponent,
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
		BacktestViewComponent,
		BacktestSubnavComponent,
		BacktestHomeComponent,
		BacktestCreateComponent,
		BacktestComingSoonComponent,

		// Reports
		AccountSummaryComponent,
		AccountReturnsComponent,
		ProfitLossComponent,
		AccountCashComponent,
		AccountValuesComponent,
		BaseComponent,

		// Trading
		IvrComponent,
		OrdersComponent,
		SubnavComponent,
		TradesComponent,
		ScreenerComponent,
		DashboardComponent,
		WatchlistComponent,
		PositionsComponent,
		MarketQuotesComponent,
		DashboardChartComponent,
		TypeaheadSymbolsComponent,
		TradeComponent,
		TradeMultiLegComponent,
		DropdownSelectComponent,
		ScreenerAddEditComponent,
		AccountHistoryComponent,
		EquityComponent,
		TradeEquityComponent,
		BrokersComponent,
		AccountComponent,
		TradingComponent,
		PersonalInfoComponent,
		BillingHistoryComponent,
		AccountDetailsComponent,
		SocialComponent,
		UpgradeComponent,
		ExpiredComponent,
		CreditCardComponent,
		CardComponent,
		ToolTipComponent,
		CustomReportsComponent,
		SpreadComponent,
		OptionComponent,
		OtherComponent,
		CentcomUsersComponent,
		CentcomCoreComponent,
		SymbolComponent,
		IronCondorComponent,
		ReverseIronCondorComponent,
		TradeOptionComponent,
		LongCallButterflyComponent,
		LongPutButterflyComponent,
		AnalyzeComponent,
	],

	imports: [
		Routing,
		FormsModule,
		BrowserModule,
		SortablejsModule,
		HttpClientModule,
		FontAwesomeModule,

		// Remember whenever we change the service worker we need to delete it in the browser. (unregister)
		environment.production ? ServiceWorkerModule.register('ngsw-worker.js') : []
	],

	providers: [
		MeService,
		QuotesService,
		AuthGuard,
		SymbolService,
		BrokerService,
		TradeGroupService,
		StateService,
		WatchlistService,
		WebsocketService,
		StatusService,
		TradeService,
		OptionsChainService,
		ScreenerService,
		ReportsService,
		AnalyzeService,
		BrokerEventsService,
		NotificationsService,
		SettingsService,
		BacktestService,
		{ provide: HTTP_INTERCEPTORS, useClass: TokenInterceptor, multi: true }
	],

	bootstrap: [AppComponent]
})

export class AppModule { }
