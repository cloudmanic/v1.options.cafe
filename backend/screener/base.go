//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

type Base struct {
	Broker      brokers.Api
	DB          models.Datastore
	ScreenFuncs map[string]func(models.Screener) ([]Result, error)
}

// Result struct
type Result struct {
	Day             types.Date      `gorm:"type:date" sql:"not null" json:"day"` // Used for backtesting. The day in the backtest
	Debit           float64         `json:"debit"`
	Credit          float64         `json:"credit"`
	MidPoint        float64         `json:"midpoint"`
	Expired         models.Date     `gorm:"type:expired" json:"expired"` // Used when all legs have the same expire
	CallPrecentAway float64         `json:"call_percent_away"`
	PutPrecentAway  float64         `json:"put_percent_away"`
	UnderlyingLast  float64         `json:"underlyng_last"`
	Legs            []models.Symbol `json:"legs"`
}

type Spread struct {
	Short types.OptionsChainItem
	Long  types.OptionsChainItem
}

type IronCondor struct {
	CallShort types.OptionsChainItem
	CallLong  types.OptionsChainItem
	PutShort  types.OptionsChainItem
	PutLong   types.OptionsChainItem
}

//
// New Screen
//
func NewScreen(db models.Datastore, broker brokers.Api) Base {

	// New screener instance
	t := Base{
		Broker: broker,
		DB:     db,
	}

	// Build screener function.
	t.ScreenFuncs = map[string]func(models.Screener) ([]Result, error){
		"iron-condor":         t.RunIronCondor,
		"short-strangle":      t.RunShortStrangle,
		"put-credit-spread":   t.RunPutCreditSpread,
		"reverse-iron-condor": t.RunReverseIronCondor,
	}

	return t
}

//
// Loop through one user's screeners and add them to our redis cache.
// This is called from our polling calls and run through our workers.
//
func PrimeScreenerCachesByUser(db models.Datastore, api brokers.Api, user models.User, broker models.Broker) error {

	// new screen instance
	t := NewScreen(db, api)

	// Get screener for this user.
	screeners := []models.Screener{}
	db.New().Where("user_id = ?", user.Id).Find(&screeners)

	for _, row := range screeners {

		// Get screen with screen items
		screen, err := db.GetScreenerByIdAndUserId(row.Id, row.UserId)

		if err != nil {
			services.Info(err)
			continue
		}

		// Run the screen based from our function map
		result, err := t.ScreenFuncs[row.Strategy](screen)

		if err != nil {
			services.Info(err)
			continue
		}

		// Store result in cache.
		cache.SetExpire("oc-screener-result-"+strconv.Itoa(int(row.Id)), (time.Minute * 5), result)

	}

	// Return happy
	return nil
}

//
// Sort results
//
func (t *Base) SortResults(result []Result) {

	// Sort the results expire in asc order.
	sort.Slice(result, func(i, j int) bool {

		// Deal with tied sorts
		if result[i].Expired.Unix() == result[j].Expired.Unix() {
			return result[i].MidPoint < result[j].MidPoint
		}

		return result[i].Expired.Unix() < result[j].Expired.Unix()
	})

}

//
// Get last quote for a symbol we are screening.
//
func (t *Base) GetQuote(smb string) (types.Quote, error) {

	// Make call to get a current quote.
	q, err := t.Broker.GetQuotes([]string{smb})

	if err != nil {
		return types.Quote{}, err
	}

	if len(q) <= 0 {
		return types.Quote{}, errors.New("No quote found.")
	}

	// return quote
	return q[0], nil
}

//
// Search filter items for a particular value.
//
func (t *Base) FindFilterItemsByKey(item string, screen models.Screener) []models.ScreenerItem {

	rt := []models.ScreenerItem{}

	for _, row := range screen.Items {

		if row.Key == item {
			rt = append(rt, row)
		}

	}

	// Return finds
	return rt
}

//
// Search filter items for a particular value. Used when we only want one value. Mostly used with the "=" operator
//
func (t *Base) FindFilterItemValue(item string, screen models.Screener) (models.ScreenerItem, error) {

	for _, row := range screen.Items {

		if row.Key == item {
			return row, nil
		}

	}

	// Not found
	return models.ScreenerItem{}, errors.New("Item not found.")
}

//
// Find a strike price that is X number of strikes below.
//
func (t *Base) FindByStrike(chain []types.OptionsChainItem, strike float64) (types.OptionsChainItem, error) {

	for _, row := range chain {

		if strike == row.Strike {
			return row, nil
		}

	}

	return types.OptionsChainItem{}, errors.New("No leg found.")
}

// ---------------- Shared Filters --------------------- //

//
// Screen final results by percent away.
//
func (t *Base) FilterPercentAwayResults(keyItem string, screen models.Screener, percentAway float64) bool {

	awayItems := t.FindFilterItemsByKey(keyItem, screen)

	// Loop over percent aways
	for _, row := range awayItems {

		// Switch based on the operator
		switch row.Operator {

		// If Percent Away > value passed in
		case "<":
			if percentAway > row.ValueNumber {
				return false
			}
			break

		// If Percent Away is less than value passed in
		case ">":
			if percentAway < row.ValueNumber {
				return false
			}
			break

		// If Percent Away is equal to values passed in
		case "=":
			if percentAway != row.ValueNumber {
				return false
			}
			break

		}

	}

	return true
}

//
// Review trade for strike percent up
//
func (t *Base) FilterStrikeByPercentUp(key string, screen models.Screener, strike float64, lastQuote float64) bool {

	var minStrike float64 = 0.00
	var widthIncrment float64 = 0.50

	// Find keys related to this filter.
	items := t.FindFilterItemsByKey(key, screen)

	// Loop through the keys. Return false if something does not match up.
	for _, row := range items {

		// Figure out the strike price that is the min we can sell.
		var tmp float64 = lastQuote + (lastQuote * (1 * (row.ValueNumber / 100)))

		fraction := tmp - math.Floor(tmp)

		if fraction >= widthIncrment {
			minStrike = (math.Floor(tmp) + widthIncrment)
		} else {
			minStrike = math.Floor(tmp)
		}

		// Switch based on the operator
		switch row.Operator {

		// Is the strike > than this value.
		case "<":
			if strike > minStrike {
				return false
			}
			break

		// Is the strike < this value
		case ">":
			if strike < minStrike {
				return false
			}
			break

		// Is the strike = this value
		case "=":
			if strike != minStrike {
				return false
			}
			break

		}

	}

	// If we made it this far it is true.
	return true
}

//
// Review trade for strike percent down
//
func (t *Base) FilterStrikeByPercentDown(key string, screen models.Screener, strike float64, lastQuote float64) bool {

	if strike > lastQuote {
		return false
	}

	var minSellStrike float64 = 0.00
	var widthIncrment float64 = 0.50

	// Find keys related to this filter.
	items := t.FindFilterItemsByKey(key, screen)

	// Loop through the keys. Return false if something does not match up.
	for _, row := range items {

		// Figure out the strike price that is the min we can sell.
		var tmp float64 = lastQuote - (lastQuote * (row.ValueNumber / 100))
		fraction := tmp - math.Floor(tmp)

		if fraction >= widthIncrment {
			minSellStrike = (math.Floor(tmp) + widthIncrment)
		} else {
			minSellStrike = math.Floor(tmp)
		}

		// Switch based on the operator
		switch row.Operator {

		// Is the strike < than this value.
		case "<":
			if strike < minSellStrike {
				return false
			}
			break

		// Is the strike > this value
		case ">":
			if strike > minSellStrike {
				return false
			}
			break

		// Is the strike = this value
		case "=":
			if strike != minSellStrike {
				return false
			}
			break

		}

	}

	// If we made it this far it is true.
	return true
}

//
// Review trade for open credit
//
func (t *Base) FilterOpenCredit(screen models.Screener, credit float64) bool {

	// Find keys related to this filter.
	items := t.FindFilterItemsByKey("open-credit", screen)

	// Loop through the keys. Return false if something does not match up.
	for _, row := range items {

		// Switch based on the operator
		switch row.Operator {

		// Is the credit > than this value.
		case "<":
			if credit > row.ValueNumber {
				return false
			}
			break

		// Is the credit < this value
		case ">":
			if credit < row.ValueNumber {
				return false
			}
			break

		// Is the credit = this value
		case "=":
			if credit != row.ValueNumber {
				return false
			}
			break

		}

	}

	// If we made it this far it is true.
	return true
}

//
// Review trade for open debit
//
func (t *Base) FilterOpenDebit(screen models.Screener, debit float64) bool {

	// Find keys related to this filter.
	items := t.FindFilterItemsByKey("open-debit", screen)

	// Loop through the keys. Return false if something does not match up.
	for _, row := range items {

		// Switch based on the operator
		switch row.Operator {

		// Is the debit > than this value.
		case "<":
			if debit > row.ValueNumber {
				return false
			}
			break

		// Is the debit < this value
		case ">":
			if debit < row.ValueNumber {
				return false
			}
			break

		// Is the debit = this value
		case "=":
			if debit != row.ValueNumber {
				return false
			}
			break

		}

	}

	// If we made it this far it is true.
	return true
}

//
// Review trade for days to expire.
//
func (t *Base) FilterDaysToExpireDaysToExpire(today time.Time, screen models.Screener, expire time.Time) bool {

	// Get days to expire.
	daysToExpire := int(today.Sub(expire).Hours()/24) * -1

	// Find keys related to this filter.
	items := t.FindFilterItemsByKey("days-to-expire", screen)

	// Loop through the keys. Return false if something does not match up.
	for _, row := range items {

		// Switch based on the operator
		switch row.Operator {

		// Is this spread to far out?
		case "<":
			if float64(daysToExpire) > row.ValueNumber {
				return false
			}
			break

		// Is this spread to far in?
		case ">":
			if float64(daysToExpire) < row.ValueNumber {
				return false
			}
			break

		// Days to expire must be equal to the value we pass in.
		case "=":
			if float64(daysToExpire) != row.ValueNumber {
				return false
			}
			break

		}

	}

	// If we made it this far it is true.
	return true
}

//
// Get possible legs
//
func (t *Base) GetPossibleVerticalSpreads(screen models.Screener, quote types.Quote, chain []types.OptionsChainItem, legType string, openType string) []Spread {

	spreads := []Spread{}

	var spreadWidth float64

	for _, row := range chain {

		// No need to pay attention to open interest of zero
		if row.OpenInterest == 0 {
			continue
		}

		// Skip strikes that are higher than our min strike. Based on percent away.
		if legType == "put" {

			if !t.FilterStrikeByPercentDown(legType+"-leg-percent-away", screen, row.Strike, quote.Last) {
				continue
			}

		} else {

			if !t.FilterStrikeByPercentUp(legType+"-leg-percent-away", screen, row.Strike, quote.Last) {
				continue
			}

		}

		// See if we have a spread width
		sw, err := t.FindFilterItemValue(legType+"-leg-width", screen)

		if err == nil {
			spreadWidth = sw.ValueNumber
		} else {
			continue
		}

		// Deal with the case of put leg
		if legType == "put" {

			// Find the strike that is x points away.
			ol, err := t.FindByStrike(chain, (row.Strike - spreadWidth))

			if err != nil {
				continue
			}

			// Add to possible to return
			if openType == "debit" {

				spreads = append(spreads, Spread{
					Short: ol,
					Long:  row,
				})

			} else {

				spreads = append(spreads, Spread{
					Short: row,
					Long:  ol,
				})

			}

		} else {

			// Find the strike that is x points away.
			ol, err := t.FindByStrike(chain, (row.Strike + spreadWidth))

			if err != nil {
				continue
			}

			// Add to possible to return
			if openType == "debit" {

				spreads = append(spreads, Spread{
					Short: ol,
					Long:  row,
				})

			} else {

				spreads = append(spreads, Spread{
					Short: row,
					Long:  ol,
				})

			}

		}

	}

	// Return happy
	return spreads
}

/* End File */
