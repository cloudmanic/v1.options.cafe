//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"crypto/rand"
	"errors"
	"html/template"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/cloudmanic/app.options.cafe/backend/library/checkmail"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	Id                 uint      `gorm:"primary_key" json:"id"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
	FirstName          string    `sql:"not null" json:"first_name"`
	LastName           string    `sql:"not null" json:"last_name"`
	Email              string    `sql:"not null" json:"email"`
	Password           string    `sql:"not null" json:"-"`
	Phone              string    `sql:"not null" json:"phone"`
	Address            string    `sql:"not null" json:"address"`
	City               string    `sql:"not null" json:"city"`
	State              string    `sql:"not null" json:"state"`
	Zip                string    `sql:"not null" json:"zip"`
	Country            string    `sql:"not null" json:"country"`
	Admin              string    `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"-"`
	Status             string    `sql:"not null;type:ENUM('Active', 'Disable', 'Delinquent', 'Expired', 'Trial');default:'Trial'" json:"status"`
	Session            Session   `json:"-"`
	Brokers            []Broker  `json:"brokers"`
	StripeCustomer     string    `sql:"not null" json:"-"`
	StripeSubscription string    `sql:"not null" json:"-"`
	GoogleSubId        string    `sql:"not null" json:"google_sub_id"`
	LastActivity       time.Time `json:"last_activity"`
	TrialExpire        time.Time `json:"-"`
	Bootstrapped       string    `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"-"` // Set to yes after we do our first import
}

type UserSubscription struct {
	Name               string    `json:"name"`
	Amount             float64   `json:"amount"`
	Status             string    `json:"status"`
	Started            time.Time `json:"started"`
	TrialStart         time.Time `json:"trial_start"`
	TrialEnd           time.Time `json:"trial_end"`
	CurrentPeriodStart time.Time `json:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end"`
	TrialDays          int       `json:"trial_days"`
	BillingInterval    string    `json:"billing_interval"`
	CardBrand          string    `json:"card_brand"`
	CardLast4          string    `json:"card_last_4"`
	CardExpMonth       int       `json:"card_exp_month"`
	CardExpYear        int       `json:"card_exp_year"`
	CouponName         string    `json:"coupon_name"`
	CouponCode         string    `json:"coupon_code"`
	CouponAmountOff    int64     `json:"coupon_amount_off"`
	CouponPercentOff   float64   `json:"coupon_percent_off"`
	CouponDuration     string    `json:"coupon_duration"`
}

type UserInvoice struct {
	Date          time.Time `json:"date"`
	Amount        float64   `json:"amount"`
	Transaction   string    `json:"transaction"`
	PaymentMethod string    `json:"payment_method"`
	InvoiceUrl    string    `json:"invoice_url"`
}

//
// Validate for this model.
//
func (a User) Validate(db Datastore, userId uint) error {

	// Return validation
	return validation.ValidateStruct(&a,

		// First Name
		validation.Field(&a.FirstName, validation.Required.Error("The first name field is required.")),

		// Last Name
		validation.Field(&a.LastName, validation.Required.Error("The last name field is required.")),

		// Email
		validation.Field(&a.Email,
			validation.Required.Error("The email field is required."),
			validation.NewStringRule(govalidator.IsEmail, "The email field must be a valid email address"),
			validation.By(func(value interface{}) error { return db.ValidateUserEmail(userId, value.(string)) }),
		),
	)
}

//
// Validate Email
//
func (t *DB) ValidateUserEmail(userId uint, email string) error {

	// Make sure this email is not already in use.
	user, _ := t.GetUserByEmail(email)

	// If we pass in the same value for email do nothing
	if (user.Id == 0) || (user.Id == userId) {
		return nil
	}

	return errors.New("Email address is already in use.")
}

//
// Delete a user and all the data that goes with that user.
//
func (t *DB) DeleteUser(user *User) error {

	// Start deleting data for this user
	t.DB.New().Where("user_id = ?", user.Id).Delete(Watchlist{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(WatchlistSymbol{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(TradeGroup{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(Settings{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(Session{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(Screener{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(Position{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(Order{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(OrderLeg{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(NotifyChannel{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(Notification{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(ForgotPassword{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(Broker{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(BrokerEvent{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(BrokerAccount{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(BalanceHistory{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(ActiveSymbol{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(Backtest{})
	t.DB.New().Where("user_id = ?", user.Id).Delete(BacktestPosition{})
	t.DB.New().Where("id = ?", user.Id).Delete(User{})

	// At Sendy make sure this user is in a canceled state
	go services.SendySubscribe("subscribers", user.Email, user.FirstName, user.LastName, "", "", "", "Yes")

	// Return happy.
	return nil
}

//
// Update user.
//
func (t *DB) UpdateUser(user *User) error {
	t.Save(user)
	return nil
}

//
// Get a user by Id.
//
func (t *DB) GetUserById(id uint) (User, error) {

	var u User

	if t.Where("id = ?", id).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Add in brokers
	t.Model(u).Related(&u.Brokers)

	// Return the user.
	return u, nil

}

//
// Verify we have default watchlist in place.
//
func (t *DB) SetDefaultWatchList(user User) {

	// Setup defaults.
	type Y struct {
		SymShort string
		SymLong  string
	}

	var m []Y
	m = append(m, Y{SymShort: "SPY", SymLong: "SPDR S&P 500"})
	m = append(m, Y{SymShort: "IWM", SymLong: "Ishares Russell 2000 Etf"})
	m = append(m, Y{SymShort: "MCD", SymLong: "McDonald's Corp"})
	m = append(m, Y{SymShort: "XLF", SymLong: "SPDR Select Sector Fund - Financial"})
	m = append(m, Y{SymShort: "AMZN", SymLong: "Amazon.com Inc"})
	m = append(m, Y{SymShort: "AAPL", SymLong: "Apple Inc."})
	m = append(m, Y{SymShort: "SBUX", SymLong: "Starbucks Corp"})
	m = append(m, Y{SymShort: "BAC", SymLong: "Bank Of America Corporation"})
	m = append(m, Y{SymShort: "HD", SymLong: "The Home Depot Inc"})
	m = append(m, Y{SymShort: "CAT", SymLong: "Caterpillar Inc"})

	// See if this user already had a watchlist
	_, err := t.GetWatchlistsByUserId(user.Id)

	// If no watchlists we create a default one with some default symbols.
	if err != nil {

		wList, err := t.CreateWatchlist(user.Id, "My Watchlist")

		if err != nil {
			services.InfoMsg(err.Error() + "(CreateNewWatchlist) Unable to create watchlist Default")
			return
		}

		for key, row := range m {

			// Add some default symbols - SPY
			symb, err := t.CreateNewSymbol(row.SymShort, row.SymLong, "Equity")

			if err != nil {
				services.InfoMsg(err.Error() + "(VerifyDefaultWatchList) Unable to create symbol " + row.SymShort)
				return
			}

			// Add lookup
			_, err2 := t.CreateWatchlistSymbol(wList, symb, user, uint(key))

			if err2 != nil {
				services.InfoMsg(err2.Error() + "(CreateNewWatchlistSymbol) Unable to create symbol " + row.SymShort + " lookup")
				return
			}

		}

	}

	return

}

//
// Get a user by Google Sub.
//
func (t *DB) GetUserByGoogleSubId(sub string) (User, error) {

	var u User

	if t.Where("google_sub_id = ?", sub).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Add in brokers
	t.Model(u).Related(&u.Brokers)

	// Return the user.
	return u, nil

}

//
// Get a user by email.
//
func (t *DB) GetUserByEmail(email string) (User, error) {

	var u User

	if t.Where("email = ?", email).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Add in brokers
	t.Model(u).Related(&u.Brokers)

	// Return the user.
	return u, nil

}

//
// Get a user by stripe customer.
//
func (t *DB) GetUserByStripeCustomer(customerId string) (User, error) {

	var u User

	if t.Where("stripe_customer = ?", customerId).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Return the user.
	return u, nil

}

//
// Return an array of all users.
//
func (t *DB) GetAllUsers() []User {

	var users []User

	t.Find(&users)

	// Add in our one to many look ups
	for i := range users {
		t.Model(users[i]).Related(&users[i].Brokers)
	}

	return users

}

//
// Return an array of all active users.
//
func (t *DB) GetAllActiveUsers() []User {

	var users []User

	t.Where("status = ?", "Active").Find(&users)

	// Add in our one to many look ups
	for i := range users {
		t.Model(users[i]).Related(&users[i].Brokers)
	}

	return users

}

//
// Return an array of all active or Trial users.
//
func (t *DB) GetAllActiveOrTrialUsers() []User {

	var users []User

	t.Where("status = ? OR status = ?", "Active", "Trial").Find(&users)

	// Add in our one to many look ups
	for i := range users {
		t.Model(users[i]).Related(&users[i].Brokers)
	}

	return users

}

//
// Login a user by ID
//
func (t *DB) LoginUserById(id uint, appId uint, userAgent string, ipAddress string) (User, error) {

	var user User

	// See if we already have this user.
	user, err := t.GetUserById(id)

	if err != nil {
		return user, errors.New("Sorry, we were unable to find our account.")
	}

	// Create a session so we get an access_token (if we passed in an appId)
	if appId > 0 {
		session, err := t.CreateSession(user.Id, appId, userAgent, ipAddress)

		if err != nil {
			services.InfoMsg(err.Error() + "LoginUserById - Unable to create session in CreateSession()")
			return User{}, err
		}

		// Add the session to the user object.
		user.Session = session
	}

	return user, nil
}

//
// Login a user in by email and password. The userAgent is a way to marking what device this
// login request came from. Same with ipAddress.
//
func (t *DB) LoginUserByEmailPass(email string, password string, appId uint, userAgent string, ipAddress string) (User, error) {

	var user User

	// See if we already have this user.
	user, err := t.GetUserByEmail(email)

	if err != nil {
		return user, errors.New("Sorry, we were unable to find our account.")
	}

	// Validate password here by comparing hashes nil means success
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return user, err
	}

	// Create a session so we get an access_token (if we passed in an appId)
	if appId > 0 {
		session, err := t.CreateSession(user.Id, appId, userAgent, ipAddress)

		if err != nil {
			services.InfoMsg(err.Error() + "LoginUserByEmailPass - Unable to create session in CreateSession()")
			return User{}, err
		}

		// Add the session to the user object.
		user.Session = session
	}

	return user, nil
}

//
// Reset a user password.
//
func (t *DB) ResetUserPassword(id uint, password string) error {

	// Get the user.
	user, err := t.GetUserById(id)

	if err != nil {
		return err
	}

	// Build the new password hash.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		services.InfoMsg(err.Error() + "ResetUserPassword - Unable to create password hash (password hash)")
		return err
	}

	// Update the database with the new password
	if err := t.Model(&user).Update("password", hash).Error; err != nil {
		services.InfoMsg(err.Error() + "ResetUserPassword - Unable update the password (mysql query)")
		return err
	}

	// Success.
	return nil

}

//
// Create a new user. - From google auth
//
func (t *DB) CreateUserFromGoogle(first string, last string, email string, subId string, appId uint, userAgent string, ipAddress string) (User, error) {

	// Lets do some validation
	if err := t.ValidateCreateUser(first, last, email, true); err != nil {
		return User{}, err
	}

	// Days to expire
	daysToExpire := helpers.StringToInt(os.Getenv("TRIAL_DAY_COUNT"))

	// Trail expire
	now := time.Now()
	tExpire := now.Add(time.Hour * 24 * time.Duration(daysToExpire))

	// Install user into the database
	var _first = template.HTMLEscapeString(first)
	var _last = template.HTMLEscapeString(last)

	// Create new user
	user := User{FirstName: _first, LastName: _last, Email: email, GoogleSubId: subId, Status: "Trial", TrialExpire: tExpire}
	t.Create(&user)

	// Log user creation.
	services.InfoMsg("CreateUser - Created a new user account (from Google Auth) - " + first + " " + last + " " + email)

	// Create a session so we get an access_token
	session, err := t.CreateSession(user.Id, appId, userAgent, ipAddress)

	if err != nil {
		services.InfoMsg(err.Error() + "CreateUser - Unable to create session in CreateSession()")
		return User{}, err
	}

	// Add the session to the user object.
	user.Session = session

	// Do post register stuff
	t.doPostUserRegisterStuff(user, ipAddress)

	// Return the user.
	return user, nil
}

//
// Create a new user.
//
func (t *DB) CreateUser(first string, last string, email string, password string, appId uint, userAgent string, ipAddress string) (User, error) {

	// Lets do some validation
	if err := t.ValidateCreateUser(first, last, email, false); err != nil {
		return User{}, err
	}

	// Make sure the password is at least 6 chars long
	if err := t.ValidatePassword(password); err != nil {
		return User{}, err
	}

	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		services.InfoMsg(err.Error() + "CreateUser - Unable to create password hash (password hash)")
		return User{}, err
	}

	// Days to expire
	daysToExpire := helpers.StringToInt(os.Getenv("TRIAL_DAY_COUNT"))

	// Trail expire
	now := time.Now()
	tExpire := now.Add(time.Hour * 24 * time.Duration(daysToExpire))

	// Install user into the database
	var _first = template.HTMLEscapeString(first)
	var _last = template.HTMLEscapeString(last)

	user := User{FirstName: _first, LastName: _last, Email: email, Password: string(hash), Status: "Trial", TrialExpire: tExpire}
	t.Create(&user)

	// Log user creation.
	services.InfoMsg("CreateUser - Created a new user account - " + first + " " + last + " " + email)

	// Create a session so we get an access_token
	session, err := t.CreateSession(user.Id, appId, userAgent, ipAddress)

	if err != nil {
		services.InfoMsg(err.Error() + "CreateUser - Unable to create session in CreateSession()")
		return User{}, err
	}

	// Add the session to the user object.
	user.Session = session

	// Do post register stuff
	t.doPostUserRegisterStuff(user, ipAddress)

	// Return the user.
	return user, nil

}

// ------------------ Stripe Functions --------------------- //

//
// Create new user with strip.
//
func (t *DB) CreateNewUserWithStripe(user User, plan string, token string, coupon string) error {

	if len(os.Getenv("STRIPE_SECRET_KEY")) > 0 {

		// First lets see if we have a customer with stripe
		custId := user.StripeCustomer

		if len(user.StripeCustomer) <= 0 {

			// Subscribe the new customer to services.
			custId2, err := services.StripeAddCustomer(user.FirstName, user.LastName, user.Email, int(user.Id))

			if err != nil {
				services.InfoMsg(err.Error() + "CreateNewUserWithStripe - Unable to create a customer account at services. - " + user.Email)
				return err
			}

			custId = custId2

		}

		// Update the user to include subscription and customer ids from strip.
		user.StripeCustomer = custId
		t.Save(&user)

		// Add the credit card to stripe
		if len(token) > 0 {
			err := t.UpdateCreditCard(user, token)

			if err != nil {
				services.Info(err)
				return err
			}
		}

		// Subscribe this user to our default Stripe plan.
		subId := user.StripeSubscription

		if len(user.StripeSubscription) <= 0 {

			subId2, err := services.StripeAddSubscription(custId, plan, coupon)

			if err != nil {
				services.InfoMsg(err.Error() + "CreateNewUserWithStripe - Unable to create a subscription at services. - " + user.Email)
				return err
			}

			subId = subId2

		}

		// Update the user to include subscription and customer ids from strip.
		user.Status = "Active"
		user.StripeSubscription = subId
		t.Save(&user)

		// Update Sendy with this new fact.
		go services.SendyUnsubscribe("trial", user.Email)
		go services.SendyUnsubscribe("expired", user.Email)
		go services.SendySubscribe("subscribers", user.Email, user.FirstName, user.LastName, "Yes", "", "", "No")

	} else {

		// Here we are doing local development or something.
		// Really we should not be doing this but sometimes we want to
		// develop stuff and not worry about strip.
		user.StripeCustomer = "na"
		user.StripeSubscription = "na"
		t.Save(&user)

	}

	// Return happy.
	return nil

}

//
// Delete user with strip.
//
func (t *DB) DeleteUserWithStripe(user User) error {

	if len(os.Getenv("STRIPE_SECRET_KEY")) > 0 {

		// Delete customer to services.
		err := services.StripeDeleteCustomer(user.StripeCustomer)

		if err != nil {
			services.InfoMsg(err.Error() + "DeleteUserWithStripe - Unable to delete a customer account at services. - " + user.Email)
			return err
		}

		// Update the user to remove ids from strip.
		user.StripeCustomer = ""
		user.StripeSubscription = ""
		t.Save(&user)

	} else {

		// Here we are doing local development or something.
		// Really we should not be doing this but sometimes we want to
		// develop stuff and not worry about strip.
		user.StripeCustomer = "na"
		user.StripeSubscription = "na"
		t.Save(&user)

	}

	// Return happy.
	return nil

}

//
// Update credit card on file. This is used if there
// is no credit card on file too.
//
func (t *DB) UpdateCreditCard(user User, stripeToken string) error {

	if len(os.Getenv("STRIPE_SECRET_KEY")) > 0 {

		// List cards in service.
		cards, err := services.StripeListAllCreditCards(user.StripeCustomer)

		if err != nil {
			services.InfoMsg(err.Error() + "AddCreditCardByToken - Unable to add a card to account. - " + user.Email)
			return err
		}

		// If we already have a card on file delete it.
		if len(cards) > 0 {
			for _, row := range cards {
				services.StripeDeleteCreditCard(user.StripeCustomer, row)
			}
		}

		// Add card to services.
		_, err2 := services.StripeAddCreditCardByToken(user.StripeCustomer, stripeToken)

		if err2 != nil {
			services.InfoMsg(err.Error() + "AddCreditCardByToken - Unable to add a card to account. - " + user.Email)
			return err
		}

	} else {
		return errors.New("No stripe key found.")
	}

	// Return happy.
	return nil

}

//
// Apply a coupon to a user's subscription
//
func (t *DB) ApplyCoupon(user User, couponCode string) error {

	if len(os.Getenv("STRIPE_SECRET_KEY")) > 0 {

		// Apply coupon to subscription.
		err := services.StripeApplyCoupon(user.StripeSubscription, couponCode)

		if err != nil {
			services.InfoMsg(err.Error() + "ApplyCoupon - Unable to apply a coupon. - " + user.Email)
			return err
		}

	} else {
		return errors.New("No stripe key found.")
	}

	// Return happy.
	return nil

}

//
// Get Stripe Subscription
//
func (t *DB) GetSubscriptionWithStripe(user User) (UserSubscription, error) {

	subscription := UserSubscription{}

	if len(os.Getenv("STRIPE_SECRET_KEY")) > 0 {

		// Get customer to services.
		cust, err := services.StripeGetCustomer(user.StripeCustomer)

		if err != nil {
			services.InfoMsg(err.Error() + "GetSubscriptionWithStripe - Unable to get customer. - " + user.Email)
			return UserSubscription{}, err
		}

		// Make sure we have at least one subscription
		if cust.Subscriptions.ListMeta.TotalCount <= 0 {
			return UserSubscription{}, errors.New("No stripe subscription found.")
		}

		// Build our internal object
		subscription.CurrentPeriodStart = time.Unix(cust.Subscriptions.Data[0].CurrentPeriodStart, 0)
		subscription.CurrentPeriodEnd = time.Unix(cust.Subscriptions.Data[0].CurrentPeriodEnd, 0)
		subscription.Name = cust.Subscriptions.Data[0].Plan.Nickname
		subscription.BillingInterval = string(cust.Subscriptions.Data[0].Plan.Interval)
		subscription.Amount = float64(cust.Subscriptions.Data[0].Plan.Amount / 100)
		subscription.TrialDays = helpers.StringToInt(os.Getenv("TRIAL_DAY_COUNT")) // we keep track our own Trail
		subscription.Status = string(cust.Subscriptions.Data[0].Status)
		subscription.Started = user.CreatedAt
		subscription.TrialStart = user.CreatedAt
		subscription.TrialEnd = user.TrialExpire

		// Do we have a credit card on file
		if cust.Sources.ListMeta.TotalCount > 0 {
			subscription.CardBrand = string(cust.Sources.Data[0].Card.Brand)
			subscription.CardLast4 = cust.Sources.Data[0].Card.Last4
			subscription.CardExpMonth = int(cust.Sources.Data[0].Card.ExpMonth)
			subscription.CardExpYear = int(cust.Sources.Data[0].Card.ExpYear)
		}

		// Do we have a coupon on file?
		if cust.Subscriptions.Data[0].Discount != nil {
			subscription.CouponCode = string(cust.Subscriptions.Data[0].Discount.Coupon.ID)
			subscription.CouponName = string(cust.Subscriptions.Data[0].Discount.Coupon.Name)
			subscription.CouponAmountOff = cust.Subscriptions.Data[0].Discount.Coupon.AmountOff
			subscription.CouponPercentOff = cust.Subscriptions.Data[0].Discount.Coupon.PercentOff
			subscription.CouponDuration = string(cust.Subscriptions.Data[0].Discount.Coupon.Duration)
		}

	} else {
		return UserSubscription{}, errors.New("No stripe key found.")
	}

	// Return happy.
	return subscription, nil

}

//
// Get invoice history from stripe
//
func (t *DB) GetInvoiceHistoryWithStripe(user User) ([]UserInvoice, error) {

	invoices := []UserInvoice{}

	// Add trial period
	invoices = append(invoices, UserInvoice{
		Date:          user.CreatedAt,
		Amount:        0,
		Transaction:   "Trial Period " + user.CreatedAt.Format("1/2/06") + " - " + user.TrialExpire.Format("1/2/06"),
		PaymentMethod: "",
		InvoiceUrl:    "",
	})

	if len(os.Getenv("STRIPE_SECRET_KEY")) > 0 {

		// Get charges by customer
		charges, err := services.StripeGetChargesByCustomer(user.StripeCustomer)

		if err != nil {
			services.InfoMsg(err.Error() + "GetInvoiceHistoryWithStripe - Unable to get charges by customer. - " + user.Email)
			return []UserInvoice{}, err
		}

		// Loop through and add invoices
		for _, row := range charges {

			// Build invoice object
			tmp := UserInvoice{
				Date:          time.Unix(row.Created, 0),
				Amount:        float64(row.Amount / 100),
				Transaction:   "Charge",
				PaymentMethod: string(row.Source.Card.Brand) + " ending " + string(row.Source.Card.Last4),
				InvoiceUrl:    "",
			}

			// is this a refund?
			if row.Refunded {
				tmp.Amount = float64(row.AmountRefunded/100) * -1
				tmp.Transaction = tmp.Transaction + " Refund"
			}

			// Add invoice information
			if row.Invoice != nil {
				inv, err := services.StripeGetInvoice(row.Invoice.ID)

				if err == nil {
					tmp.InvoiceUrl = inv.InvoicePDF

					// Get the date range
					start := time.Unix(inv.Lines.Data[0].Period.Start, 0).Format("1/2/06")
					end := time.Unix(inv.Lines.Data[0].Period.End, 0).Format("1/2/06")
					tmp.Transaction = "Subscription " + start + " - " + end

				}

			}

			invoices = append(invoices, tmp)

		}

	}

	// Return happy.
	return invoices, nil
}

// ------------------ Helper Functions --------------------- //

//
// Do post user register stuff.
//
func (t *DB) doPostUserRegisterStuff(user User, ipAddress string) {

	// Set default watchlists
	t.SetDefaultWatchList(user)

	// Set required active symbols
	t.CreateActiveSymbol(user.Id, "$DJI")
	t.CreateActiveSymbol(user.Id, "COMP")
	t.CreateActiveSymbol(user.Id, "VIX")
	t.CreateActiveSymbol(user.Id, "SPX")

	// Subscribe new user to mailing lists.
	go services.SendySubscribe("trial", user.Email, user.FirstName, user.LastName, "", "", ipAddress, "No")
	go services.SendySubscribe("no-brokers", user.Email, user.FirstName, user.LastName, "No", "", ipAddress, "No")
	go services.SendySubscribe("subscribers", user.Email, user.FirstName, user.LastName, "No", "", ipAddress, "No")

	// Tell slack about this.
	go services.SlackNotify("#events", "New Options Cafe User Account : "+user.Email)

	// Add our welcome notice.
	title := "Welcome to Options Cafe Beta"
	message := `Thanks for giving us a try! You are using a beta release. During our beta release we will be focusing on <a href="https://www.investopedia.com/university/optionspreadstrategies/optionspreads2.asp" target="_blank">vertical spread</a> options strategies. More strategies will roll out over the course of the next month. If you have any issues or want to give us feedback please email us at <a href="mailto:help@options.cafe">help@options.cafe</a>.`

	n := Notification{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserId:      user.Id,
		Status:      "pending",
		Channel:     "in-app",
		Uri:         "dashboard-notice",
		Title:       title,
		LongMessage: message,
		SentTime:    time.Now(),
		Expires:     time.Now().AddDate(0, 1, 0), // 1 month
	}

	// Store in DB
	t.DB.New().Save(&n)
}

//
// Validate a login user action.
//
func (t *DB) ValidateUserLogin(email string, password string) error {

	// Make sure the password is at least 6 chars long
	if len(password) < 6 {
		return errors.New("The password filed must be at least 6 characters long.")
	}

	// Lets validate the email address
	if err := t.ValidateEmailAddress(email); err != nil {
		return err
	}

	// See if we already have this user.
	_, err := t.GetUserByEmail(email)

	if err != nil {
		return errors.New("Sorry, we were unable to find our account.")
	}

	// Return happy.
	return nil
}

//
// Validate a create user action. We do not always get a first name and last name from google.
// so we make the validation optional with them.
//
func (t *DB) ValidateCreateUser(first string, last string, email string, googleAuth bool) error {

	// Are first and last name fields empty
	if (!googleAuth) && (len(first) == 0) && (len(last) == 0) {
		return errors.New("First name and last name fields are required.")
	}

	// Are first name empty
	if (!googleAuth) && len(first) == 0 {
		return errors.New("First name field is required.")
	}

	// Are last name empty
	if (!googleAuth) && len(last) == 0 {
		return errors.New("Last name field is required.")
	}

	// Lets validate the email address
	if err := t.ValidateEmailAddress(email); err != nil {
		return err
	}

	// See if we already have this user.
	_, err := t.GetUserByEmail(email)

	if err == nil {
		return errors.New("Looks like you already have an account.")
	}

	// Return happy.
	return nil
}

//
// Validate password.
//
func (t *DB) ValidatePassword(password string) error {

	// Make sure the password is at least 6 chars long
	if len(password) < 6 {
		return errors.New("The password filed must be at least 6 characters long.")
	}

	// Return happy.
	return nil

}

//
// Validate an email address
//
func (t *DB) ValidateEmailAddress(email string) error {

	// Check length
	if len(email) == 0 {
		return errors.New("Email address field is required.")
	}

	// Check format
	if err := checkmail.ValidateFormat(email); err != nil {
		return errors.New("Email address is not a valid format.")
	}

	// Return happy.
	return nil

}

//
// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
//
func (t *DB) GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

	bytes, err := t.GenerateRandomBytes(n)

	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes), nil
}

//
// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
//
func (t *DB) GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

/* End File */
