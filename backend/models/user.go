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

	"github.com/cloudmanic/app.options.cafe/backend/library/checkmail"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                 uint      `gorm:"primary_key" json:"id"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
	FirstName          string    `sql:"not null" json:"first_name"`
	LastName           string    `sql:"not null" json:"last_name"`
	Email              string    `sql:"not null" json:"email"`
	Password           string    `sql:"not null" json:"-"`
	Admin              string    `sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"-"`
	Status             string    `sql:"not null;type:ENUM('Active', 'Disable');default:'Active'" json:"-"`
	Session            Session   `json:"-"`
	Brokers            []Broker  `json:"brokers"`
	StripeCustomer     string    `sql:"not null" json:"-"`
	StripeSubscription string    `sql:"not null" json:"-"`
	GoogleSubId        string    `sql:"not null" json:"-"`
	LastActivity       time.Time `json:"last_activity"`
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
			services.Error(err, "LoginUserById - Unable to create session in CreateSession()")
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
			services.Error(err, "LoginUserByEmailPass - Unable to create session in CreateSession()")
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
		services.Error(err, "ResetUserPassword - Unable to create password hash (password hash)")
		return err
	}

	// Update the database with the new password
	if err := t.Model(&user).Update("password", hash).Error; err != nil {
		services.Error(err, "ResetUserPassword - Unable update the password (mysql query)")
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

	// Install user into the database
	var _first = template.HTMLEscapeString(first)
	var _last = template.HTMLEscapeString(last)

	// Create new user
	user := User{FirstName: _first, LastName: _last, Email: email, GoogleSubId: subId}
	t.Create(&user)

	// Log user creation.
	services.Info("CreateUser - Created a new user account (from Google Auth) - " + first + " " + last + " " + email)

	// Create a session so we get an access_token
	session, err := t.CreateSession(user.Id, appId, userAgent, ipAddress)

	if err != nil {
		services.Error(err, "CreateUser - Unable to create session in CreateSession()")
		return User{}, err
	}

	// Add the session to the user object.
	user.Session = session

	// Talk to stripe and setup the account.
	err = t.CreateNewUserWithStripe(user)

	if err != nil {
		return User{}, err
	}

	// Do post register stuff
	t.doPostUserRegisterStuff(user)

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
		services.Error(err, "CreateUser - Unable to create password hash (password hash)")
		return User{}, err
	}

	// Install user into the database
	var _first = template.HTMLEscapeString(first)
	var _last = template.HTMLEscapeString(last)

	user := User{FirstName: _first, LastName: _last, Email: email, Password: string(hash)}
	t.Create(&user)

	// Log user creation.
	services.Info("CreateUser - Created a new user account - " + first + " " + last + " " + email)

	// Create a session so we get an access_token
	session, err := t.CreateSession(user.Id, appId, userAgent, ipAddress)

	if err != nil {
		services.Error(err, "CreateUser - Unable to create session in CreateSession()")
		return User{}, err
	}

	// Add the session to the user object.
	user.Session = session

	// Talk to stripe and setup the account.
	err = t.CreateNewUserWithStripe(user)

	if err != nil {
		return User{}, err
	}

	// Do post register stuff
	t.doPostUserRegisterStuff(user)

	// Return the user.
	return user, nil

}

// ------------------ Helper Functions --------------------- //

//
// Do post user register stuff.
//
func (t *DB) doPostUserRegisterStuff(user User) {

	// Subscribe new user to mailing lists.
	go services.SendySubscribe("no-brokers", user.Email, user.FirstName, user.LastName)
	go services.SendySubscribe("subscribers", user.Email, user.FirstName, user.LastName)

	// Tell slack about this.
	go services.SlackNotify("#events", "New Options Cafe User Account : "+user.Email)

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
// Create new user with strip.
//
func (t *DB) CreateNewUserWithStripe(user User) error {

	if len(os.Getenv("STRIPE_SECRET_KEY")) > 0 {

		// Subscribe the new customer to services.
		custId, err := services.AddCustomer(user.FirstName, user.LastName, user.Email, int(user.Id))

		if err != nil {
			services.Error(err, "CreateNewUserWithStripe - Unable to create a customer account at services. - "+user.Email)
			return err
		}

		// Subscribe this user to our default Stripe plan.
		subId, err := services.AddSubscription(custId, os.Getenv("STRIPE_DEFAULT_PLAN"))

		if err != nil {
			services.Error(err, "CreateNewUserWithStripe - Unable to create a subscription at services. - "+user.Email)
			return err
		}

		// Update the user to include subscription and customer ids from strip.
		user.StripeCustomer = custId
		user.StripeSubscription = subId
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
