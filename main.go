package main

import "fmt"
import "gopkg.in/gomail.v2"

func main() {
	fmt.Println("hello world")
	msg := gomail.NewMessage()
	msg.SetHeader("From", "hello@pinmyblogs.com")
	msg.SetHeader("To", "punitinani1@gmail.com")
	msg.SetHeader("Subject", "Test Mail pinmyblogs.com")
	msg.SetBody("text/html", "<b>This is the body of the mail</b>")

	n := gomail.NewDialer(
		"smtpout.secureserver.net",
		465,
		"hello@pinmyblogs.com",
		"password",
	)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}
}
