package models

import (
  "time"
  "errors"
  "crypto/rand"
  "golang.org/x/crypto/bcrypt"
  "app.options.cafe/backend/library/services"
)

type User struct {
  Id uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  FirstName string `sql:"not null"`
  LastName string `sql:"not null"`
  Email string `sql:"not null"`
  Password string `sql:"not null"`
  AccessToken string `sql:"not null"`
  Status string `sql:"not null;type:ENUM('Active', 'Disable');default:'Active'"`
  Brokers []Broker
}

//
// Get a user by email.
//
func (t * DB) GetUserByEmail(email string) (User, error) {
 
  var u User
  
  if t.Connection.Where("email = ?", email).First(&u).RecordNotFound() {
    return u, errors.New("Record not found")
  }
  
  // Return the user.
  return u, nil
  
}     
      
// 
// Return an array of all users.
//
func (t * DB) GetAllUsers() ([]User) {
  
  var users []User
  
  t.Connection.Find(&users)
  
  // Add in our one to many look ups
  for i, _ := range users {
    t.Connection.Model(users[i]).Related(&users[i].Brokers)     
  }  
  
  return users
  
}

// 
// Return an array of all active users.
//
func (t * DB) GetAllActiveUsers() ([]User) {
  
  var users []User
  
  t.Connection.Where("status = ?", "Active").Find(&users)
  
  // Add in our one to many look ups
  for i, _ := range users {
    t.Connection.Model(users[i]).Related(&users[i].Brokers)     
  }  
  
  return users
  
}

//
// Create a new user.
//
func (t * DB) CreateUser(first string, last string, email string, password string) (User, error) {
  
  // Create an access token.
  access_token, err := GenerateRandomString(50)

	if err != nil {
    services.Error(err, "CreateUser - Unable to create random string (access_token)")
    return User{}, err 
	}  
	
  // Generate "hash" to store from user password
  hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  
  if err != nil {
    services.Error(err, "CreateUser - Unable to create password hash (password hash)")
    return User{}, err    
  }
  
  // Install user into the database
  user := User{FirstName: first, LastName: last, Email: email, Password: string(hash), AccessToken: access_token}
  t.Connection.Create(&user)
  
  // Log user creation.
  services.Log("CreateUser - Created a new user account - " + first + " " + last + " " + email)
 
  // Return the user.
  return user, nil
  
}

// ------------------ Helper Functions --------------------- //

//
// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
//
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	
	bytes, err := GenerateRandomBytes(n)
	
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
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	
	_, err := rand.Read(b)
	
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
      
/* End File */