package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail() {
	// Get html
	var body bytes.Buffer
	t, err := template.ParseFiles(os.Getenv("MAIL_TEMPLATE"))
	t.Execute(&body, struct{ Name string }{Name: "VeMines"})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send with GoMail
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("GOOGLE_USRENAME"))
	m.SetHeader("To", os.Getenv("GOOGLE_USRENAME"))
	m.SetHeader("Subject", "Hello! Subject Here")
	m.SetBody("text/html", body.String())
	m.Attach("./assets/uploads/image.png")

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("GOOGLE_USRENAME"), os.Getenv("GOOGLE_APP_PASSWORD"))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
