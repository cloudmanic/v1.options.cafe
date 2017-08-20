package models

import (
  "time"
  "errors"
  "app.options.cafe/backend/library/email"
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
func (t * DB) DoResetPassword(user_email string, ip string) error {

  // Make sure this is a real email address
  user, err := t.GetUserByEmail(user_email)
  
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
  
  // Send email to user asking them to come to the site and reset the password.
  err = email.Send(user.Email, "Options Cafe Reset Password Request", GetForgotPasswordStepOneEmailHtml())

  if err != nil { 
    return err
  }

  // Everything went as planned.
  return nil
  
}

// ------------------- Template Emails ------------------------- //

func GetForgotPasswordStepOneEmailHtml() string {
  
  return string(`
    <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
            "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
    <html>
    
    </head>
    
    <body>
    <p>
        Hello {{.Name}}
        <a href="{{.URL}}">Confirm email address</a>
    </p>
        
    </body>
    
    </html>  
  `)
  
}
 
/* End File */