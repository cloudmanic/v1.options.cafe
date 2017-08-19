package models

import (
  "time"
  "errors"
  "app.options.cafe/backend/library/services"
)

type ForgotPassword struct {
  Id uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  UserId uint `sql:"not null;index:UserId"` 
  Token string `sql:"not null"`
  IpAddress string `sql:"not null"`
}

//
// Reset the user's password and send an email telling them next steps.
//
func (t * DB) DoResetPassword(email string, ip string) error {

  // Make sure this is a real email address
  user, err := t.GetUserByEmail(email)
  
  if err != nil {  
    return errors.New("Sorry, we were unable to find our account.")
  }

  // Generate "hash" to store for the reset token
  hash, err := GenerateRandomString(30)
  
  if err != nil {
    services.Error(err, "DoResetPassword - Unable to create hash (GenerateRandomString)")
    return err    
  }
  
  // Store the new reset password hash.
  rsp := ForgotPassword{UserId: user.Id, Token: hash, IpAddress: ip }
  t.Connection.Create(&rsp)
  
  // Log user creation.
  services.Log("DoResetPassword - Reset password token for " + user.Email)   

  // Everything went as planned.
  return nil
  
}
 
/* End File */