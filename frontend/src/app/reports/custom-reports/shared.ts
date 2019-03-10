//
// Date: 10/20/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

export class Shared {
	// List timeframes
	static readonly TimeFrames: TimeFrame[] = [
		{ Value: "1-year", Name: "Last 1 Year" },
		{ Value: "2-year", Name: "Last 2 Year" },
		{ Value: "3-year", Name: "Last 3 Year" },
		{ Value: "4-year", Name: "Last 4 Year" },
		{ Value: "5-year", Name: "Last 5 Year" },
		{ Value: "10-year", Name: "Last 10 Year" },
		{ Value: "ytd", Name: "Year to Date" },
		{ Value: "custom", Name: "Custom Dates" }
	]

	// List custom report types
	static readonly ReportTypes: ReportType[] = [
		{ Name: "Profit & Loss", Key: "profit-loss", Route: "/reports/custom/profit-loss", Query: false },
		{ Name: "Account Cash", Key: "account-cash", Route: "/reports/custom/account-cash", Query: false },
		{ Name: "Account Value", Key: "account-values", Route: "/reports/custom/account-values", Query: false },
		{ Name: "Account Returns", Key: "account-returns", Route: "/reports/custom/account-returns", Query: false },
		{ Name: "Cumulative Earnings", Key: "profit-loss-cumulative", Route: "/reports/custom/profit-loss", Query: true },
	]
}

export interface TimeFrame {
	Value: string;
	Name: string;
}

export interface ReportType {
	Key: string;
	Name: string;
	Route: string;
	Query: boolean;
}
