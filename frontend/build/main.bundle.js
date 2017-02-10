webpackJsonp([0,3],{

/***/ 302:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return Watchlist; });
var Watchlist = (function () {
    function Watchlist(id, name, items) {
        this.id = id;
        this.name = name;
        this.items = items;
    }
    return Watchlist;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/watchlist.js.map

/***/ },

/***/ 344:
/***/ function(module, exports) {

function webpackEmptyContext(req) {
	throw new Error("Cannot find module '" + req + "'.");
}
webpackEmptyContext.keys = function() { return []; };
webpackEmptyContext.resolve = webpackEmptyContext;
module.exports = webpackEmptyContext;
webpackEmptyContext.id = 344;


/***/ },

/***/ 345:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__polyfills_ts__ = __webpack_require__(470);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__polyfills_ts___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0__polyfills_ts__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__ = __webpack_require__(432);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__environments_environment__ = __webpack_require__(469);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__app_app_module__ = __webpack_require__(454);





if (__WEBPACK_IMPORTED_MODULE_3__environments_environment__["a" /* environment */].production) {
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_2__angular_core__["_37" /* enableProdMode */])();
}
__webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__["a" /* platformBrowserDynamic */])().bootstrapModule(__WEBPACK_IMPORTED_MODULE_4__app_app_module__["a" /* AppModule */]);
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/main.js.map

/***/ },

/***/ 452:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_broker_service__ = __webpack_require__(51);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return AccountsComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var AccountsComponent = (function () {
    //
    // Constructor....
    //
    function AccountsComponent(broker, changeDetect) {
        this.broker = broker;
        this.changeDetect = changeDetect;
    }
    //
    // OnInit....
    //
    AccountsComponent.prototype.ngOnInit = function () {
        var _this = this;
        // Subscribe to data updates from the broker - Market Status
        this.broker.userProfilePushData.subscribe(function (data) {
            _this.userProfile = data;
            // Do we have an account already? Always have to reset the selected one when we get new account data.
            if ((!_this.selectedAccount) && (_this.userProfile.Accounts.length)) {
                _this.selectedAccount = _this.userProfile.Accounts[0];
                _this.broker.setActiveAccountId(_this.selectedAccount.AccountNumber);
            }
            else {
                for (var i = 0; i < _this.userProfile.Accounts.length; i++) {
                    if (_this.userProfile.Accounts[i].AccountNumber == _this.selectedAccount.AccountNumber) {
                        _this.selectedAccount = _this.userProfile.Accounts[i];
                    }
                }
            }
            _this.changeDetect.detectChanges();
        });
    };
    //
    // On account change.
    //
    AccountsComponent.prototype.onAccountChange = function () {
        this.broker.setActiveAccountId(this.selectedAccount.AccountNumber);
    };
    AccountsComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["G" /* Component */])({
            selector: 'oc-accounts',
            template: __webpack_require__(623)
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_broker_service__["a" /* BrokerService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__services_broker_service__["a" /* BrokerService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */]) === 'function' && _b) || Object])
    ], AccountsComponent);
    return AccountsComponent;
    var _a, _b;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/accounts.component.js.map

/***/ },

/***/ 453:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return AppComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};

var AppComponent = (function () {
    //
    // Constructor....
    //
    function AppComponent() {
    }
    //
    // On Init
    //
    AppComponent.prototype.ngOnInit = function () {
    };
    AppComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["G" /* Component */])({
            selector: 'oc-root',
            template: __webpack_require__(624),
        }), 
        __metadata('design:paramtypes', [])
    ], AppComponent);
    return AppComponent;
}());
/* End File */
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/app.component.js.map

/***/ },

/***/ 454:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__ = __webpack_require__(193);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_forms__ = __webpack_require__(422);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_http__ = __webpack_require__(428);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__app_component__ = __webpack_require__(453);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__services_quote_service__ = __webpack_require__(96);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__services_broker_service__ = __webpack_require__(51);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__layout_header_component__ = __webpack_require__(465);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8__layout_footer_component__ = __webpack_require__(464);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9__accounts_accounts_component__ = __webpack_require__(452);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10__watchlist_watchlist_component__ = __webpack_require__(468);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11__orders_orders_component__ = __webpack_require__(466);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12__dashboard_dashboard_component__ = __webpack_require__(463);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13__sidebar_sidebar_component__ = __webpack_require__(467);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return AppModule; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};














var AppModule = (function () {
    function AppModule() {
    }
    AppModule = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_core__["I" /* NgModule */])({
            declarations: [
                __WEBPACK_IMPORTED_MODULE_4__app_component__["a" /* AppComponent */],
                __WEBPACK_IMPORTED_MODULE_7__layout_header_component__["a" /* HeaderComponent */],
                __WEBPACK_IMPORTED_MODULE_8__layout_footer_component__["a" /* FooterComponent */],
                __WEBPACK_IMPORTED_MODULE_9__accounts_accounts_component__["a" /* AccountsComponent */],
                __WEBPACK_IMPORTED_MODULE_10__watchlist_watchlist_component__["a" /* WatchlistComponent */],
                __WEBPACK_IMPORTED_MODULE_11__orders_orders_component__["a" /* OrdersComponent */],
                __WEBPACK_IMPORTED_MODULE_12__dashboard_dashboard_component__["a" /* DashboardComponent */],
                __WEBPACK_IMPORTED_MODULE_13__sidebar_sidebar_component__["a" /* SidebarComponent */]
            ],
            imports: [
                __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__["b" /* BrowserModule */],
                __WEBPACK_IMPORTED_MODULE_2__angular_forms__["a" /* FormsModule */],
                __WEBPACK_IMPORTED_MODULE_3__angular_http__["a" /* HttpModule */]
            ],
            providers: [__WEBPACK_IMPORTED_MODULE_6__services_broker_service__["a" /* BrokerService */], __WEBPACK_IMPORTED_MODULE_5__services_quote_service__["a" /* QuoteService */]],
            bootstrap: [__WEBPACK_IMPORTED_MODULE_4__app_component__["a" /* AppComponent */]]
        }), 
        __metadata('design:paramtypes', [])
    ], AppModule);
    return AppModule;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/app.module.js.map

/***/ },

/***/ 455:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return Balance; });
var Balance = (function () {
    function Balance(AccountNumber, AccountValue, TotalCash, OptionBuyingPower, StockBuyingPower) {
        this.AccountNumber = AccountNumber;
        this.AccountValue = AccountValue;
        this.TotalCash = TotalCash;
        this.OptionBuyingPower = OptionBuyingPower;
        this.StockBuyingPower = StockBuyingPower;
    }
    return Balance;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/balance.js.map

/***/ },

/***/ 456:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return BrokerAccounts; });
var BrokerAccounts = (function () {
    function BrokerAccounts(AccountNumber, Classification, DayTrader, OptionLevel, Status, Type) {
        this.AccountNumber = AccountNumber;
        this.Classification = Classification;
        this.DayTrader = DayTrader;
        this.OptionLevel = OptionLevel;
        this.Status = Status;
        this.Type = Type;
    }
    return BrokerAccounts;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/broker-accounts.js.map

/***/ },

/***/ 457:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return MarketQuote; });
var MarketQuote = (function () {
    //
    // Constructor...
    //
    function MarketQuote(last, open, prev_close, symbol, description) {
        this.last = last;
        this.open = open;
        this.prev_close = prev_close;
        this.symbol = symbol;
        this.description = description;
    }
    //
    // Return the daily change
    //
    MarketQuote.prototype.dailyChange = function () {
        return (this.last - this.prev_close).toFixed(2);
    };
    //
    // Return the percent change
    //
    MarketQuote.prototype.percentChange = function () {
        // Do we have data yet?
        if (this.prev_close <= 0) {
            return 0;
        }
        return parseFloat((((this.last - this.prev_close) / this.prev_close) * 100).toFixed(2));
    };
    //
    // Return the class color
    //
    MarketQuote.prototype.classColor = function () {
        var change = this.percentChange();
        if (change > 0) {
            return 'green';
        }
        else if (change < 0) {
            return 'red';
        }
        return '';
    };
    return MarketQuote;
}());
/* End File */ 
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/market-quote.js.map

/***/ },

/***/ 458:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return MarketStatus; });
var MarketStatus = (function () {
    function MarketStatus(description, state) {
        this.description = description;
        this.state = state;
    }
    return MarketStatus;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/market-status.js.map

/***/ },

/***/ 459:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return OrderLeg; });
var OrderLeg = (function () {
    function OrderLeg(Type, Symbol, OptionSymbol, Side, Quantity, Status, Duration, AvgFillPrice, ExecQuantity, LastFillPrice, LastFillQuantity, RemainingQuantity, CreateDate, TransactionDate) {
        this.Type = Type;
        this.Symbol = Symbol;
        this.OptionSymbol = OptionSymbol;
        this.Side = Side;
        this.Quantity = Quantity;
        this.Status = Status;
        this.Duration = Duration;
        this.AvgFillPrice = AvgFillPrice;
        this.ExecQuantity = ExecQuantity;
        this.LastFillPrice = LastFillPrice;
        this.LastFillQuantity = LastFillQuantity;
        this.RemainingQuantity = RemainingQuantity;
        this.CreateDate = CreateDate;
        this.TransactionDate = TransactionDate;
    }
    return OrderLeg;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/order-leg.js.map

/***/ },

/***/ 460:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return Order; });
var Order = (function () {
    function Order(Id, AccountId, AvgFillPrice, Class, CreateDate, Duration, ExecQuantity, LastFillPrice, LastFillQuantity, NumLegs, Price, Quantity, RemainingQuantity, Side, Status, Symbol, TransactionDate, Type, Legs) {
        this.Id = Id;
        this.AccountId = AccountId;
        this.AvgFillPrice = AvgFillPrice;
        this.Class = Class;
        this.CreateDate = CreateDate;
        this.Duration = Duration;
        this.ExecQuantity = ExecQuantity;
        this.LastFillPrice = LastFillPrice;
        this.LastFillQuantity = LastFillQuantity;
        this.NumLegs = NumLegs;
        this.Price = Price;
        this.Quantity = Quantity;
        this.RemainingQuantity = RemainingQuantity;
        this.Side = Side;
        this.Status = Status;
        this.Symbol = Symbol;
        this.TransactionDate = TransactionDate;
        this.Type = Type;
        this.Legs = Legs;
    }
    return Order;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/order.js.map

/***/ },

/***/ 461:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return UserProfile; });
var UserProfile = (function () {
    function UserProfile(Id, Name, Accounts) {
        this.Id = Id;
        this.Name = Name;
        this.Accounts = Accounts;
    }
    return UserProfile;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/user-profile.js.map

/***/ },

/***/ 462:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return WatchlistItems; });
var WatchlistItems = (function () {
    function WatchlistItems(id, symbol) {
        this.id = id;
        this.symbol = symbol;
    }
    return WatchlistItems;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/watchlist-items.js.map

/***/ },

/***/ 463:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_quote_service__ = __webpack_require__(96);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__services_broker_service__ = __webpack_require__(51);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return DashboardComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var DashboardComponent = (function () {
    function DashboardComponent(quotesService, broker, changeDetect) {
        this.quotesService = quotesService;
        this.broker = broker;
        this.changeDetect = changeDetect;
        this.ws_reconnecting = false;
    }
    //
    // OnInit....
    //
    DashboardComponent.prototype.ngOnInit = function () {
        var _this = this;
        // Subscribe to when we are reconnecting to a websocket - Core
        this.broker.wsReconnecting.subscribe(function (data) {
            _this.ws_reconnecting = data;
            _this.changeDetect.detectChanges();
        });
        // Subscribe to when we are reconnecting to a websocket - Quotes
        this.quotesService.wsReconnecting.subscribe(function (data) {
            _this.ws_reconnecting = data;
            _this.changeDetect.detectChanges();
        });
    };
    DashboardComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["G" /* Component */])({
            selector: 'oc-dashboard',
            template: __webpack_require__(625)
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_quote_service__["a" /* QuoteService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__services_quote_service__["a" /* QuoteService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__services_broker_service__["a" /* BrokerService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__services_broker_service__["a" /* BrokerService */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */]) === 'function' && _c) || Object])
    ], DashboardComponent);
    return DashboardComponent;
    var _a, _b, _c;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/dashboard.component.js.map

/***/ },

/***/ 464:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_broker_service__ = __webpack_require__(51);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return FooterComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var FooterComponent = (function () {
    //
    // Constructor....
    //
    function FooterComponent(broker, changeDetect) {
        this.broker = broker;
        this.changeDetect = changeDetect;
    }
    //
    // OnInit....
    //
    FooterComponent.prototype.ngOnInit = function () {
        var _this = this;
        // Subscribe to data updates from the broker - Market Status
        this.broker.marketStatusPushData.subscribe(function (data) {
            _this.marketStatus = data;
            _this.changeDetect.detectChanges();
        });
    };
    FooterComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["G" /* Component */])({
            selector: 'oc-footer',
            template: __webpack_require__(626)
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_broker_service__["a" /* BrokerService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__services_broker_service__["a" /* BrokerService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */]) === 'function' && _b) || Object])
    ], FooterComponent);
    return FooterComponent;
    var _a, _b;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/footer.component.js.map

/***/ },

/***/ 465:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_broker_service__ = __webpack_require__(51);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__services_quote_service__ = __webpack_require__(96);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return HeaderComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var HeaderComponent = (function () {
    //
    // Constructor....
    //
    function HeaderComponent(broker, quotesService, changeDetect) {
        this.broker = broker;
        this.quotesService = quotesService;
        this.changeDetect = changeDetect;
        this.quotes = {};
    }
    //
    // OnInit....
    //
    HeaderComponent.prototype.ngOnInit = function () {
        var _this = this;
        // Subscribe to data updates from the broker - Market Status
        this.broker.userProfilePushData.subscribe(function (data) {
            _this.userProfile = data;
            _this.changeDetect.detectChanges();
        });
        // Subscribe to data updates from the quotes - Market Quotes
        this.quotesService.marketQuotePushData.subscribe(function (data) {
            _this.quotes[data.symbol] = data;
            _this.changeDetect.detectChanges();
        });
    };
    HeaderComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["G" /* Component */])({
            selector: 'oc-header',
            template: __webpack_require__(627)
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_broker_service__["a" /* BrokerService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__services_broker_service__["a" /* BrokerService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__services_quote_service__["a" /* QuoteService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__services_quote_service__["a" /* QuoteService */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */]) === 'function' && _c) || Object])
    ], HeaderComponent);
    return HeaderComponent;
    var _a, _b, _c;
}());
/* End File */ 
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/header.component.js.map

/***/ },

/***/ 466:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_quote_service__ = __webpack_require__(96);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__services_broker_service__ = __webpack_require__(51);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return OrdersComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var OrdersComponent = (function () {
    //
    // Constructor....
    //
    function OrdersComponent(quotesService, broker, changeDetect) {
        this.quotesService = quotesService;
        this.broker = broker;
        this.changeDetect = changeDetect;
        this.quotes = {};
        this.activeAccount = "";
    }
    //
    // OnInit....
    //
    OrdersComponent.prototype.ngOnInit = function () {
        var _this = this;
        // Set the active account.
        this.activeAccount = this.broker.activeAccount;
        // Subscribe to data updates from the broker - Orders
        this.broker.ordersPushData.subscribe(function (data) {
            //console.log(data);
            var rt = [];
            // Filter - We only one the accounts that are active.
            for (var i = 0; i < data.length; i++) {
                if (data[i].AccountId == _this.activeAccount) {
                    rt.push(data[i]);
                }
            }
            // Set order data
            _this.orders = rt;
            _this.changeDetect.detectChanges();
        });
        // Subscribe to when the active account changes
        this.broker.activeAccountPushData.subscribe(function (data) {
            _this.activeAccount = data;
            _this.orders = [];
            _this.changeDetect.detectChanges();
        });
        // Subscribe to data updates from the quotes - Market Quotes
        this.quotesService.marketQuotePushData.subscribe(function (data) {
            _this.quotes[data.symbol] = data;
            _this.changeDetect.detectChanges();
        });
    };
    OrdersComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["G" /* Component */])({
            selector: 'oc-orders',
            template: __webpack_require__(628)
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_quote_service__["a" /* QuoteService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__services_quote_service__["a" /* QuoteService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__services_broker_service__["a" /* BrokerService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__services_broker_service__["a" /* BrokerService */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */]) === 'function' && _c) || Object])
    ], OrdersComponent);
    return OrdersComponent;
    var _a, _b, _c;
}());
/* End File */ 
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/orders.component.js.map

/***/ },

/***/ 467:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_broker_service__ = __webpack_require__(51);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return SidebarComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var SidebarComponent = (function () {
    //
    // Constructor....
    //
    function SidebarComponent(broker, changeDetect) {
        this.broker = broker;
        this.changeDetect = changeDetect;
    }
    //
    // OnInit....
    //
    SidebarComponent.prototype.ngOnInit = function () {
        var _this = this;
        // Subscribe to data updates from the broker - User Profile
        this.broker.userProfilePushData.subscribe(function (data) {
            _this.userProfile = data;
            _this.changeDetect.detectChanges();
        });
        // Subscribe to data updates from the broker - Market Status
        this.broker.marketStatusPushData.subscribe(function (data) {
            _this.marketStatus = data;
            _this.changeDetect.detectChanges();
        });
        // Subscribe to data updates from the broker - Balances
        this.broker.balancesPushData.subscribe(function (data) {
            for (var i = 0; i < data.length; i++) {
                if (data[i].AccountNumber == _this.broker.activeAccount) {
                    _this.balance = data[i];
                }
            }
            _this.changeDetect.detectChanges();
        });
    };
    SidebarComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["G" /* Component */])({
            selector: 'oc-sidebar',
            template: __webpack_require__(629)
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_broker_service__["a" /* BrokerService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__services_broker_service__["a" /* BrokerService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */]) === 'function' && _b) || Object])
    ], SidebarComponent);
    return SidebarComponent;
    var _a, _b;
}());
/* End File */ 
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/sidebar.component.js.map

/***/ },

/***/ 468:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_quote_service__ = __webpack_require__(96);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__services_broker_service__ = __webpack_require__(51);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__contracts_watchlist__ = __webpack_require__(302);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return WatchlistComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};




var WatchlistComponent = (function () {
    //
    // Constructor....
    //
    function WatchlistComponent(quotesService, broker, changeDetect) {
        this.quotesService = quotesService;
        this.broker = broker;
        this.changeDetect = changeDetect;
        this.quotes = {};
        this.watchlist = new __WEBPACK_IMPORTED_MODULE_3__contracts_watchlist__["a" /* Watchlist */]('', '', []);
    }
    //
    // OnInit....
    //
    WatchlistComponent.prototype.ngOnInit = function () {
        var _this = this;
        // Subscribe to data updates from the broker - Watchlist
        this.broker.watchlistPushData.subscribe(function (data) {
            _this.watchlist = data;
            _this.changeDetect.detectChanges();
        });
        // Subscribe to data updates from the quotes - Market Quotes
        this.quotesService.marketQuotePushData.subscribe(function (data) {
            _this.quotes[data.symbol] = data;
            _this.changeDetect.detectChanges();
        });
    };
    WatchlistComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["G" /* Component */])({
            selector: 'oc-watchlist',
            template: __webpack_require__(630)
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_quote_service__["a" /* QuoteService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__services_quote_service__["a" /* QuoteService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__services_broker_service__["a" /* BrokerService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__services_broker_service__["a" /* BrokerService */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_0__angular_core__["i" /* ChangeDetectorRef */]) === 'function' && _c) || Object])
    ], WatchlistComponent);
    return WatchlistComponent;
    var _a, _b, _c;
}());
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/watchlist.component.js.map

/***/ },

/***/ 469:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return environment; });
// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `angular-cli.json`.
var environment = {
    production: false
};
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/environment.js.map

/***/ },

/***/ 470:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol__ = __webpack_require__(484);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_core_js_es6_object__ = __webpack_require__(477);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_core_js_es6_object___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_core_js_es6_object__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_core_js_es6_function__ = __webpack_require__(473);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_core_js_es6_function___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_core_js_es6_function__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int__ = __webpack_require__(479);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float__ = __webpack_require__(478);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_core_js_es6_number__ = __webpack_require__(476);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_core_js_es6_number___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_5_core_js_es6_number__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_core_js_es6_math__ = __webpack_require__(475);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_core_js_es6_math___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_6_core_js_es6_math__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_core_js_es6_string__ = __webpack_require__(483);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_core_js_es6_string___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_7_core_js_es6_string__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_core_js_es6_date__ = __webpack_require__(472);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_core_js_es6_date___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_8_core_js_es6_date__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9_core_js_es6_array__ = __webpack_require__(471);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9_core_js_es6_array___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_9_core_js_es6_array__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp__ = __webpack_require__(481);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11_core_js_es6_map__ = __webpack_require__(474);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11_core_js_es6_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_11_core_js_es6_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12_core_js_es6_set__ = __webpack_require__(482);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12_core_js_es6_set___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_12_core_js_es6_set__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect__ = __webpack_require__(480);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect__ = __webpack_require__(485);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone__ = __webpack_require__(642);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone__);
















//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/polyfills.js.map

/***/ },

/***/ 51:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__contracts_order__ = __webpack_require__(460);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__contracts_balance__ = __webpack_require__(455);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__contracts_order_leg__ = __webpack_require__(459);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__contracts_watchlist__ = __webpack_require__(302);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__contracts_watchlist_items__ = __webpack_require__(462);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__contracts_market_status__ = __webpack_require__(458);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__contracts_user_profile__ = __webpack_require__(461);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8__contracts_broker_accounts__ = __webpack_require__(456);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return BrokerService; });









var BrokerService = (function () {
    //
    // Construct!!
    //
    function BrokerService() {
        this.deviceId = "";
        this.activeAccount = "";
        // Data objects
        this.marketStatus = new __WEBPACK_IMPORTED_MODULE_6__contracts_market_status__["a" /* MarketStatus */]('', '');
        this.userProfile = new __WEBPACK_IMPORTED_MODULE_7__contracts_user_profile__["a" /* UserProfile */]('', '', []);
        // Websocket Stuff
        this.ws = null;
        this.heartbeat = null;
        this.missed_heartbeats = 0;
        // Emitters - Pushers
        this.wsReconnecting = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["_20" /* EventEmitter */]();
        this.ordersPushData = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["_20" /* EventEmitter */]();
        this.balancesPushData = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["_20" /* EventEmitter */]();
        this.userProfilePushData = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["_20" /* EventEmitter */]();
        this.marketStatusPushData = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["_20" /* EventEmitter */]();
        this.watchlistPushData = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["_20" /* EventEmitter */]();
        this.activeAccountPushData = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["_20" /* EventEmitter */]();
        // Set the device id
        var clientJs = new ClientJS();
        this.deviceId = clientJs.getFingerprint();
        // Setup standard websocket connection.
        this.setupWebSocket();
    }
    //
    // Do User Profile Refresh
    //
    BrokerService.prototype.doUserProfileRefresh = function (data) {
        // Clear accounts array.
        this.userProfile.Accounts = [];
        // Setup the array of accounts.
        for (var i in data.Accounts) {
            this.userProfile.Accounts.push(new __WEBPACK_IMPORTED_MODULE_8__contracts_broker_accounts__["a" /* BrokerAccounts */](data.Accounts[i].AccountNumber, data.Accounts[i].Classification, data.Accounts[i].DayTrader, data.Accounts[i].OptionLevel, data.Accounts[i].Status, data.Accounts[i].Type));
        }
        this.userProfile.Id = data.Id;
        this.userProfile.Name = data.Name;
        this.userProfilePushData.emit(this.userProfile);
    };
    //
    // Do Market Status Refresh
    //
    BrokerService.prototype.doMarketStatusRefresh = function (data) {
        this.marketStatus.state = data.State;
        this.marketStatus.description = data.Description;
        this.marketStatusPushData.emit(this.marketStatus);
    };
    //
    // Do Balances Refresh
    //
    BrokerService.prototype.doBalancesRefresh = function (data) {
        var balances = [];
        for (var i = 0; i < data.length; i++) {
            balances.push(new __WEBPACK_IMPORTED_MODULE_2__contracts_balance__["a" /* Balance */](data[i].AccountNumber, data[i].AccountValue, data[i].TotalCash, data[i].OptionBuyingPower, data[i].StockBuyingPower));
        }
        this.balancesPushData.emit(balances);
    };
    //
    // Do watchlist Refresh
    //
    BrokerService.prototype.doWatchListRefresh = function (data) {
        // We only care about the default watchlist
        if (data.Id != 'default') {
            return false;
        }
        var ws = new __WEBPACK_IMPORTED_MODULE_4__contracts_watchlist__["a" /* Watchlist */](data.Id, data.Name, []);
        for (var i in data.Symbols) {
            ws.items.push(new __WEBPACK_IMPORTED_MODULE_5__contracts_watchlist_items__["a" /* WatchlistItems */](data.Symbols[i].id, data.Symbols[i].symbol));
        }
        this.watchlistPushData.emit(ws);
    };
    //
    // Do refresh orders
    //
    BrokerService.prototype.doOrdersRefresh = function (data) {
        var orders = [];
        for (var i = 0; i < data.length; i++) {
            // Add in the legs
            var legs = [];
            if (data[i].NumLegs > 0) {
                for (var k = 0; k < data[i].Legs.length; k++) {
                    legs.push(new __WEBPACK_IMPORTED_MODULE_3__contracts_order_leg__["a" /* OrderLeg */](data[i].Legs[k].Type, data[i].Legs[k].Symbol, data[i].Legs[k].OptionSymbol, data[i].Legs[k].Side, data[i].Legs[k].Quantity, data[i].Legs[k].Status, data[i].Legs[k].Duration, data[i].Legs[k].AvgFillPrice, data[i].Legs[k].ExecQuantity, data[i].Legs[k].LastFillPrice, data[i].Legs[k].LastFillQuantity, data[i].Legs[k].RemainingQuantity, data[i].Legs[k].CreateDate, data[i].Legs[k].TransactionDate));
                }
            }
            // Push the order on
            orders.push(new __WEBPACK_IMPORTED_MODULE_1__contracts_order__["a" /* Order */](data[i].Id, data[i].AccountId, data[i].AvgFillPrice, data[i].Class, data[i].CreateDate, data[i].Duration, data[i].ExecQuantity, data[i].LastFillPrice, data[i].LastFillQuantity, data[i].NumLegs, data[i].Price, data[i].Quantity, data[i].RemainingQuantity, data[i].Side, data[i].Status, data[i].Symbol, data[i].TransactionDate, data[i].Type, legs));
        }
        this.ordersPushData.emit(orders);
    };
    // ------------------------ Push Data Back To Backend --------------------- //
    //
    // Request the backend sends all data again. (often do this on state change or page change)
    //
    BrokerService.prototype.requestAllData = function () {
        this.ws.send(JSON.stringify({ type: 'refresh-all-data', data: {} }));
    };
    //
    // Set the active account id.
    //
    BrokerService.prototype.setActiveAccountId = function (account_id) {
        this.activeAccount = account_id;
        this.activeAccountPushData.emit(account_id);
        this.requestAllData();
    };
    // ------------------------ Websocket Stuff --------------------- //
    //
    // Setup normal data websocket connection.
    //
    BrokerService.prototype.setupWebSocket = function () {
        var _this = this;
        // Setup websocket
        this.ws = new WebSocket(ws_server + '/ws/core');
        // Websocket sent data to us.
        this.ws.onmessage = function (e) {
            var msg = JSON.parse(e.data);
            // Is this a pong to our ping or some other return.
            if (msg.type == 'pong') {
                _this.missed_heartbeats--;
            }
            else {
                var msg_data = JSON.parse(msg.data);
                // Send quote to angular component
                switch (msg.type) {
                    // User Profile refresh
                    case 'UserProfile:refresh':
                        _this.doUserProfileRefresh(msg_data);
                        break;
                    // Market Status refresh
                    case 'MarketStatus:refresh':
                        _this.doMarketStatusRefresh(msg_data);
                        break;
                    // Watchlist refresh
                    case 'Watchlist:refresh':
                        _this.doWatchListRefresh(msg_data);
                        break;
                    // Order refresh
                    case 'Orders:refresh':
                        _this.doOrdersRefresh(msg_data);
                        break;
                    // Balances refresh
                    case 'Balances:refresh':
                        _this.doBalancesRefresh(msg_data);
                        break;
                }
            }
        };
        // On Websocket open
        this.ws.onopen = function (e) {
            // Send Access Token (Give a few moments to get started)
            setTimeout(function () {
                _this.ws.send(JSON.stringify({
                    type: 'set-access-token',
                    data: { access_token: localStorage.getItem('access_token'), device_id: _this.deviceId }
                }));
            }, 1000);
            // Tell the UI we are connected
            _this.wsReconnecting.emit(false);
            // Setup the connection heartbeat
            if (_this.heartbeat === null) {
                _this.missed_heartbeats = 0;
                _this.heartbeat = setInterval(function () {
                    try {
                        _this.missed_heartbeats++;
                        if (_this.missed_heartbeats >= 5) {
                            throw new Error('Too many missed heartbeats.');
                        }
                        _this.ws.send(JSON.stringify({ type: 'ping' }));
                    }
                    catch (e) {
                        _this.wsReconnecting.emit(true);
                        clearInterval(_this.heartbeat);
                        _this.heartbeat = null;
                        console.warn("Closing connection. Reason: " + e.message);
                        _this.ws.close();
                    }
                }, 5000);
            }
            else {
                clearInterval(_this.heartbeat);
            }
        };
        // On Close
        this.ws.onclose = function () {
            // Kill Ping heartbeat.
            clearInterval(_this.heartbeat);
            _this.heartbeat = null;
            _this.ws = null;
            // Try to reconnect
            _this.wsReconnecting.emit(true);
            setTimeout(function () { _this.setupWebSocket(); }, 3 * 1000);
        };
    };
    return BrokerService;
}());
/* End File */
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/broker.service.js.map

/***/ },

/***/ 623:
/***/ function(module, exports) {

module.exports = "<div class=\"select\">\n  <select [(ngModel)]=\"selectedAccount\" (ngModelChange)=\"onAccountChange()\">\n    <option [ngValue]=\"row\" *ngFor=\"let row of userProfile?.Accounts\">{{ row.AccountNumber }}</option>\n  </select>\n</div>"

/***/ },

/***/ 624:
/***/ function(module, exports) {

module.exports = "<oc-header></oc-header>\n\n<div class=\"main\">\n  \n  <oc-sidebar></oc-sidebar>\n  \n  <oc-dashboard></oc-dashboard>\n\n</div>\n\n<oc-footer></oc-footer>"

/***/ },

/***/ 625:
/***/ function(module, exports) {

module.exports = "<div class=\"content\">\n\t\n\t<div class=\"content-head\">\n\n\t\t<div class=\"tabs\">\n\t\t\t<ul class=\"nav nav-tabs\" role=\"tablist\">\n\t\t\t\t<li role=\"presentation\">\n\t\t\t\t\t<a href=\"\" aria-controls=\"tab1\" role=\"tab\" data-toggle=\"tab\">&nbsp;</a>\n\t\t\t\t</li>\n\t\t\t\t\n\t\t\t\t<li role=\"presentation\">\n\t\t\t\t\t<a href=\"\" aria-controls=\"tab2\" role=\"tab\" data-toggle=\"tab\">&nbsp;</a>\n\t\t\t\t</li>\n\t\t\t\t\n\t\t\t\t<li role=\"presentation\">\n\t\t\t\t\t<a href=\"\" aria-controls=\"tab3\" role=\"tab\" data-toggle=\"tab\">&nbsp;</a>\n\t\t\t\t</li>\n\t\t\t\t\n\t\t\t\t<li role=\"presentation\" class=\"active\">\n\t\t\t\t\t<a href=\"\" aria-controls=\"tab4\" role=\"tab\" data-toggle=\"tab\">Dashboard</a>\n\t\t\t\t</li>\n\t\t\t\t\n\t\t\t\t<li role=\"presentation\">\n\t\t\t\t\t<a href=\"\" aria-controls=\"tab5\" role=\"tab\" data-toggle=\"tab\">Screener</a>\n\t\t\t\t</li>\n\t\t\t\t\n\t\t\t\t<li role=\"presentation\">\n\t\t\t\t\t<a href=\"\" aria-controls=\"tab6\" role=\"tab\" data-toggle=\"tab\">Backtest</a>\n\t\t\t\t</li>\n\n\t\t\t\t<li role=\"presentation\">\n\t\t\t\t\t<a href=\"\" aria-controls=\"tab7\" role=\"tab\" data-toggle=\"tab\">&nbsp;</a>\n\t\t\t\t</li>\n\t\t\t\t\n\t\t\t\t<li role=\"presentation\">\n\t\t\t\t\t<a href=\"\" aria-controls=\"tab8\" role=\"tab\" data-toggle=\"tab\">&nbsp;</a>\n\t\t\t\t</li>\n\t\t\t\t\n\t\t\t\t<li role=\"presentation\">\n\t\t\t\t\t<a href=\"\" aria-controls=\"tab9\" role=\"tab\" data-toggle=\"tab\">&nbsp;</a>\n\t\t\t\t</li>\n\t\t\t</ul>\n\t\t</div><!-- /.tabs -->\n\t</div><!-- /.content-head -->\n\t\n  <div class=\"content-body\">\n    <div class=\"tabs\">\n\n      <div class=\"zone-content section-styleguide dashboard\">\n        \n        <div [hidden]=\"! ws_reconnecting\" class=\"bg-warning text-white\" style=\"margin-bottom: 10px;\">Reconnecting to the server...</div>\n      \t\n      \t<div class=\"row\">\n          <oc-watchlist></oc-watchlist>\n          \n          <div class=\"col-md-9\">\n            \n            <oc-orders></oc-orders>\n            \n          </div>\n          \n      \t</div>\n      \t\n      </div>\n    \n    </div>\n  </div>\n\n</div>"

/***/ },

/***/ 626:
/***/ function(module, exports) {

module.exports = ""

/***/ },

/***/ 627:
/***/ function(module, exports) {

module.exports = "<header class=\"header\">\n  <div class=\"container-fluid\">\n  \t<div class=\"header-body\">\n    \t<div class=\"logo-cont\">\n    \t  <img src=\"assets/css/images/stockpeer_logo.jpg\" width=\"40\" height=\"40\" />\n    \t</div>\n    \t\n  \t\t<div class=\"header-content\">\n        <div class=\"row\">\n          \n          <div class=\"col-md-3 text-center\">\n    \t\t\t  <strong>Dow</strong> {{ quotes['$DJI']?.last | currency:'USD':true:'1.2-2' }} <br />\n            <span class=\"{{ quotes['$DJI']?.classColor() }}\">{{ quotes['$DJI']?.dailyChange() | currency:'USD':true:'1.2-2' }} ({{ quotes['$DJI']?.percentChange() | number:'1.2-2' }}%)</span>\n    \t\t\t</div>\n    \t\t\t\n    \t\t\t<div class=\"col-md-3  text-center\">\n    \t\t\t  <strong>S&P 500</strong> {{ quotes['SPX']?.last | currency:'USD':true:'1.2-2' }}<br />\n    \t\t\t <span class=\"{{ quotes['SPX']?.classColor() }}\">{{ quotes['SPX']?.dailyChange() | currency:'USD':true:'1.2-2' }} ({{ quotes['SPX']?.percentChange() | number:'1.2-2' }}%)</span>\n    \t\t\t</div>\n\n    \t\t\t<div class=\"col-md-3  text-center\">\n    \t\t\t  <strong>Nasdaq</strong> {{ quotes['COMP']?.last | currency:'USD':true:'1.2-2' }}<br />\n            <span class=\"{{ quotes['COMP']?.classColor() }}\">{{ quotes['COMP']?.dailyChange() | currency:'USD':true:'1.2-2' }} ({{ quotes['COMP']?.percentChange() | number:'1.2-2' }}%)</span>\n    \t\t\t</div>\n\n    \t\t\t<div class=\"col-md-3 text-center\">\n    \t\t\t  <strong>VIX</strong> {{ quotes['VIX']?.last | currency:'USD':true:'1.2-2' }}<br />\n            <span class=\"{{ quotes['VIX']?.classColor() }}\">{{ quotes['VIX']?.dailyChange() | currency:'USD':true:'1.2-2' }} ({{ quotes['VIX']?.percentChange() | number:'1.2-2' }}%)</span>\n    \t\t\t</div>\n\t    \t\t\t\n  \t\t\t</div>\n  \t\t\t  \t\t\t\n  \t\t</div><!-- /.header-content -->  \t\t\n  \t\t\n  \t</div><!-- /.header-body -->\n  </div><!-- /.container-fluid -->\n</header><!-- /.header -->"

/***/ },

/***/ 628:
/***/ function(module, exports) {

module.exports = "<div class=\"panel-group\">\n  <div class=\"panel panel-default\">\n  \t<div class=\"panel-heading\" role=\"tab\" id=\"headingOne\">\n  \t\t<h4 class=\"panel-title\">\n  \t\t\t<a role=\"button\" data-toggle=\"collapse\" href=\"#collapseOrders\">\n  \t\t\t\tOrders\n  \t\t\t</a>\n  \t\t</h4>\n  \t</div>\n  \t\n  \t<div id=\"collapseOrders\" class=\"panel-collapse collapse in\" role=\"tabpanel\" aria-labelledby=\"headingOne\">\n  \t\t<div class=\"panel-body\">\n\n        <table class=\"table\">\n          <thead>\n            <tr>\n            <th>Symbols</th>\n            <th class=\"text-center\">Last</th>\n            <th class=\"text-center\">Qty</th>              \n            <th class=\"text-center\">Duration</th>       \n            <th class=\"text-center\">Type</th>\n            <th class=\"text-center\">Price</th> \n            <th class=\"text-center\">Filled</th>                       \n            <th class=\"text-center\">Status</th>        \n            </tr>\n          </thead>\n        \n          <tbody>\n\n            <tr *ngFor=\"let row of orders\">\n\n              <td [hidden]=\"row.Class != 'equity'\">{{ quotes[row.Symbol]?.description }}</td>\n              <td [hidden]=\"row.Class != 'equity'\" class=\"text-center\">{{ quotes[row.Symbol]?.last | currency:'USD':true:'1.2-2' }}</td>              \n              <td [hidden]=\"row.Class != 'equity'\" class=\"text-center\">{{ row.Quantity }}</td>\n\n\n              <td [hidden]=\"row.Class != 'multileg'\">\n                <p *ngFor=\"let row2 of row.Legs\">{{ quotes[row2.OptionSymbol]?.description }}</p>                \n              </td>\n              <td [hidden]=\"row.Class != 'multileg'\" class=\"text-center\">\n                <p *ngFor=\"let row2 of row.Legs\">{{ quotes[row2.OptionSymbol]?.last | currency:'USD':true:'1.2-2' }}</p>                \n              </td>               \n              <td [hidden]=\"row.Class != 'multileg'\" class=\"text-center\">\n                <p *ngFor=\"let row2 of row.Legs\">{{ row2.Quantity }}</p>                \n              </td>  \n\n              \n              <td [hidden]=\"row.Class != 'option'\">{{ quotes[row.Symbol]?.description }}</td>\n              <td [hidden]=\"row.Class != 'option'\" class=\"text-center\">{{ quotes[row.Symbol]?.last | currency:'USD':true:'1.2-2' }}</td> \n              <td [hidden]=\"row.Class != 'option'\" class=\"text-center\">{{ row.Quantity }}</td>\n                            \n  \n              <td class=\"text-center\">{{ row.Duration }}</td>\n              <td class=\"text-center\">{{ row.Type }}</td>\n\n              <td class=\"text-center\" [hidden]=\"row.Type != 'debit'\">{{ row.Price | currency:'USD':true:'1.2-2' }}</td>\n              <td class=\"text-center\" [hidden]=\"row.Type != 'credit'\">{{ row.Price | currency:'USD':true:'1.2-2' }}</td>\n              <td class=\"text-center\" [hidden]=\"row.Type != 'stop'\">{{ row.StopPrice | currency:'USD':true:'1.2-2' }}</td>\n              <td class=\"text-center\" [hidden]=\"row.Type != 'market'\">---</td>              \n\n              <td class=\"text-center\" [hidden]=\"row.Status != 'filled'\">{{ row.AvgFillPrice | currency:'USD':true:'1.2-2' }}</td>\n              <td class=\"text-center\" [hidden]=\"row.Status != 'open'\">---</td>\n              <td class=\"text-center\" [hidden]=\"row.Status != 'canceled'\">---</td>\n              <td class=\"text-center\" [hidden]=\"row.Status != 'expired'\">---</td>\n\n              <td class=\"text-center\">{{ row.Status }}</td>\n                                             \n            </tr>            \n            \n          </tbody>\n        </table>\n\n  \t\t</div>\n  \t</div>\n  </div>\n</div>"

/***/ },

/***/ 629:
/***/ function(module, exports) {

module.exports = "<div class=\"sidebar\">\n\t<ul class=\"widgets\">\n\t\t<li class=\"widget widget-main\">\n\t\t\t<div class=\"widget-head\">\n\t\t\t\t<h5 class=\"widget-title\">Main</h5>\n\t\t\t</div>\n\t\t\t\n\t\t\t<div class=\"widget-body\">\n\t\t\t\t<nav class=\"nav-utilities\">\n\t\t\t\t\t<ul>\n\t\t\t\t\t\t<li>\n\t\t\t\t\t\t\t<a href=\"#\">\n\t\t\t\t\t\t\t\t<span>\n\t\t\t\t\t\t\t\t\t<i class=\"ico-indicator\"></i>\n\t\t\t\t\t\t\t\t</span>\n\t\t\t\t\t\t\t\tTrading\n\t\t\t\t\t\t\t</a>\n\t\t\t\t\t\t</li>\n\t\t\t\t\t\t\n\t\t\t\t\t\t<li>\n\t\t\t\t\t\t\t<a href=\"#\">\n\t\t\t\t\t\t\t\t<span>\n\t\t\t\t\t\t\t\t\t<i class=\"ico-report\"></i>\n\t\t\t\t\t\t\t\t</span>\n\n\t\t\t\t\t\t\t\tReports\n\t\t\t\t\t\t\t</a>\n\t\t\t\t\t\t</li>\n\t\t\t\t\t</ul>\n\t\t\t\t</nav><!-- /.nav-utilities -->\n\t\t\t</div><!-- /.widget-body -->\n\t\t</li><!-- /.widget -->\n\t\t\n\t\t<li class=\"widget widget-account\">\n\t\t\t<div class=\"widget-head\">\n\t\t\t\t<h5 class=\"widget-title\">Account</h5>\n\t\t\t</div>\n\t\t\t\n\t\t\t<div class=\"widget-body\">\n  \t\t\t\n\t\t\t\t<nav class=\"nav-utilities nav-utilities-secondary\">\n\t\t\t\t\t<ul>\n\t\t\t\t\t\t<li>\n  \t\t        <oc-accounts></oc-accounts>\n\t\t\t\t\t\t</li>\n\t\t\t\t\t</ul>\n\t\t\t\t</nav>\n\t\t\t\t\n\t\t\t</div>\n\t\t\t\n\t\t\t<br />\n\t\t\t\n\t\t\t<div class=\"widget-body\">\n  \t\t\t<strong>Total Account Value</strong> <br />  \t\t\t\n  \t\t\t<span>{{ balance?.AccountValue | currency:'USD':true:'1.2-2' }}</span>\n\t\t\t</div>\n\t\t\t\n      <br />\n\t\t\t<div class=\"widget-body\">\n  \t\t\t<strong>Total Cash</strong> <br />  \t\t\t\n  \t\t\t<span>{{ balance?.TotalCash | currency:'USD':true:'1.2-2' }}</span>\n\t\t\t</div>\n\t\t\t\t\t\t\n\t\t\t<br />\n\t\t\t\n\t\t\t<div class=\"widget-body\">\n  \t\t\t<strong>Option Buying Power</strong> <br />  \t\t\t\n  \t\t\t<span>{{ balance?.OptionBuyingPower | currency:'USD':true:'1.2-2' }}</span>\n\t\t\t</div>\n\n      <br />\n\t\t\t<div class=\"widget-body\">\n  \t\t\t<strong>Stock Buying Power</strong> <br />  \t\t\t\n  \t\t\t<span>{{ balance?.StockBuyingPower | currency:'USD':true:'1.2-2' }}</span>\n\t\t\t</div>\n\t\t\t\t\t\t\n\t\t</li>\n\t\t\n\t\t<li class=\"widget\">\t\n\t\t\t<div class=\"widget-head\">\n\t\t\t\t<h5 class=\"widget-title\">Market</h5>\n\t\t\t</div>\t\t\n\t\t\t\n\t\t\t<div class=\"widget-body\">\n  \t\t\t<strong>US Stock Market</strong> <br />  \t\t\t\n  \t\t\t<span [ngSwitch]=\"marketStatus?.state\">\n  \t\t\t  <span *ngSwitchCase=\"'open'\">Open</span>\n  \t\t\t  <span *ngSwitchDefault>Closed</span>\n  \t\t\t</span>\n\t\t\t</div>\t\t\t\t\t\t\n\t\t\n\t\t</li>\t\t\n\t\t\t\n\t</ul><!-- /.widgets -->\n</div><!-- /.sidebar -->"

/***/ },

/***/ 630:
/***/ function(module, exports) {

module.exports = "<div class=\"panel-group col-md-3 watchlist\">\n  <div class=\"panel panel-default\">\n  \t<div class=\"panel-heading\" role=\"tab\" id=\"headingOne\">\n  \t\t<h4 class=\"panel-title\">\n  \t\t\t<a role=\"button\" data-toggle=\"collapse\" href=\"#collapseWatchlist\">\n  \t\t\t\tWatchlist\n  \t\t\t</a>\n  \t\t</h4>\n  \t</div>\n  \t\n  \t<div id=\"collapseWatchlist\" class=\"panel-collapse collapse in\" role=\"tabpanel\" aria-labelledby=\"headingOne\">\n  \t\t<div class=\"panel-body\">\n\n        <div class=\"list-group\">\n        \n          <a href=\"\" class=\"list-group-item\" *ngFor=\"let row of watchlist.items\">\n            <strong>{{ row.symbol }}</strong> <br /> \n            ${{ quotes[row.symbol]?.last | number:'1.2-2' }} \n            <span class=\"{{ quotes[row.symbol]?.classColor() }}\">\n              {{ quotes[row.symbol]?.dailyChange() | currency:'USD':true:'1.2-2' }}\n              <span>({{ quotes[row.symbol]?.percentChange() | number:'1.2-2' }}%)</span>\n            </span> \n            <br />\n            {{ quotes[row.symbol]?.description }}\n          </a>\n                   \n        </div>\n\n  \t\t</div>\n  \t</div>\n  </div>\n</div>"

/***/ },

/***/ 643:
/***/ function(module, exports, __webpack_require__) {

module.exports = __webpack_require__(345);


/***/ },

/***/ 96:
/***/ function(module, exports, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__contracts_market_quote__ = __webpack_require__(457);
/* harmony export (binding) */ __webpack_require__.d(exports, "a", function() { return QuoteService; });


var QuoteService = (function () {
    //
    // Construct!!
    //
    function QuoteService() {
        this.deviceId = "";
        this.quotes = {};
        // Websocket Stuff
        this.ws = null;
        this.heartbeat = null;
        this.missed_heartbeats = 0;
        // Emitters
        this.wsReconnecting = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["_20" /* EventEmitter */]();
        this.marketQuotePushData = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["_20" /* EventEmitter */]();
        // Set the device id
        var clientJs = new ClientJS();
        this.deviceId = clientJs.getFingerprint();
        // Setup standard websocket connection.
        this.setupWebSocket();
    }
    // ------------------------ Websocket Stuff --------------------- //
    //
    // Setup normal data websocket connection.
    //
    QuoteService.prototype.setupWebSocket = function () {
        var _this = this;
        // Setup websocket
        this.ws = new WebSocket(ws_server + '/ws/quotes');
        // Websocket sent data to us.
        this.ws.onmessage = function (e) {
            var msg = JSON.parse(e.data);
            // Is this a pong to our ping or some other return.
            if (msg.type == 'pong') {
                _this.missed_heartbeats--;
            }
            else {
                // Send quote to angular component
                switch (msg.type) {
                    // Real-time market quote
                    case 'trade':
                        // Have we seen this quote before?
                        if (typeof _this.quotes[msg.symbol] == "undefined") {
                            _this.quotes[msg.symbol] = new __WEBPACK_IMPORTED_MODULE_1__contracts_market_quote__["a" /* MarketQuote */](msg.last, 0, 0, msg.symbol, '');
                        }
                        else {
                            _this.quotes[msg.symbol].last = msg.last;
                        }
                        _this.marketQuotePushData.emit(_this.quotes[msg.symbol]);
                        break;
                    // DetailedQuotes refresh
                    case 'DetailedQuotes:refresh':
                        var msg_data = JSON.parse(msg.data);
                        // Have we seen this quote before?
                        if (typeof _this.quotes[msg_data.Symbol] == "undefined") {
                            _this.quotes[msg_data.Symbol] = new __WEBPACK_IMPORTED_MODULE_1__contracts_market_quote__["a" /* MarketQuote */](msg_data.Last, msg_data.Open, msg_data.PrevClose, msg_data.Symbol, msg_data.Description);
                        }
                        else {
                            _this.quotes[msg_data.Symbol].last = msg_data.Last;
                            _this.quotes[msg_data.Symbol].open = msg_data.Open;
                            _this.quotes[msg_data.Symbol].prev_close = msg_data.PrevClose;
                            _this.quotes[msg_data.Symbol].description = msg_data.Description;
                        }
                        _this.marketQuotePushData.emit(_this.quotes[msg_data.Symbol]);
                        break;
                }
            }
        };
        // On Websocket open
        this.ws.onopen = function (e) {
            // Send Access Token (Give a few moments to get started)
            setTimeout(function () {
                _this.ws.send(JSON.stringify({
                    type: 'set-access-token',
                    data: { access_token: localStorage.getItem('access_token'), device_id: _this.deviceId }
                }));
            }, 1000);
            // Tell the UI we are connected
            _this.wsReconnecting.emit(false);
            // Setup the connection heartbeat
            if (_this.heartbeat === null) {
                _this.missed_heartbeats = 0;
                _this.heartbeat = setInterval(function () {
                    try {
                        _this.missed_heartbeats++;
                        if (_this.missed_heartbeats >= 5) {
                            throw new Error('Too many missed heartbeats (quotes).');
                        }
                        _this.ws.send(JSON.stringify({ type: 'ping' }));
                    }
                    catch (e) {
                        _this.wsReconnecting.emit(true);
                        clearInterval(_this.heartbeat);
                        _this.heartbeat = null;
                        console.warn("Closing connection (quotes). Reason: " + e.message);
                        _this.ws.close();
                    }
                }, 5000);
            }
            else {
                clearInterval(_this.heartbeat);
            }
        };
        // On Close
        this.ws.onclose = function () {
            // Kill Ping heartbeat.
            clearInterval(_this.heartbeat);
            _this.heartbeat = null;
            _this.ws = null;
            // Try to reconnect
            _this.wsReconnecting.emit(true);
            setTimeout(function () { _this.setupWebSocket(); }, 3 * 1000);
        };
    };
    return QuoteService;
}());
/* End File */
//# sourceMappingURL=/Users/spicer/Development/app.options.cafe/frontend/src/quote.service.js.map

/***/ }

},[643]);
//# sourceMappingURL=main.bundle.map