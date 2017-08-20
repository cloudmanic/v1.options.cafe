package email

import (
  "os"
  "strconv"
  "gopkg.in/gomail.v2"
  "app.options.cafe/backend/library/services"
)

//
// Pass in everything we need to send an email and we send it.
// If we have a SMTP in our configs we use that if not we use
// Mailgun's library for sending mail. 
//
func Send(to string, subject string, html string, text string) error {
    
  // Setup the email to send.
  m := gomail.NewMessage()
  m.SetHeader("From", "help@options.cafe")
  m.SetHeader("To", to)
  m.SetHeader("Bcc", "bcc@options.cafe")
  m.SetHeader("Subject", subject)
  m.SetBody("text/html", html)
  m.AddAlternative("text/plain", text)

  // Configure the port to be an int. 
  port, _ := strconv.ParseInt(os.Getenv("MAIL_PORT"), 10, 64)
  
  // Make a SMTP connection
  d := gomail.NewDialer(os.Getenv("MAIL_HOST"), 
                        int(port), 
                        os.Getenv("MAIL_USERNAME"), 
                        os.Getenv("MAIL_PASSWORD"))

  // Send Da Email
  if err := d.DialAndSend(m); err != nil {
    services.Error(err, "library/email/Send/Send() - Unable to send email.") 
    return err
  }  
  
  // Everything went well!
  return nil
  
}

