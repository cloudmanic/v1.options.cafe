package email

import (
  "os"
  "bytes"
  "strconv"
  "html/template"
  "gopkg.in/gomail.v2"
  "app.options.cafe/backend/library/services"
)

//
// Pass in everything we need to send an email and we send it.
// If we have a SMTP in our configs we use that if not we use
// Mailgun's library for sending mail. 
//
func Send(to string, subject string, html string) error {
  
	templateData := struct {
		Name string
		URL  string
	}{
		Name: "Dhanush",
		URL:  "http://geektrust.in",
	}  
  
  
  // Build the html email we are going to send.
  tpl := template.New(subject)
  t, err := tpl.Parse(html)
	
	if err != nil {
    return err
	}
  
	buf := new(bytes.Buffer)
	
	if err = t.Execute(buf, templateData); err != nil {
		return err
	}	 
  
  // Setup the email to send.
  m := gomail.NewMessage()
  m.SetHeader("From", "help@options.cafe")
  m.SetHeader("To", to)
  m.SetHeader("Bcc", "bcc@options.cafe")
  m.SetHeader("Subject", subject)
  m.SetBody("text/html", buf.String())

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

