package main

import (
	"fmt"
	// "net/smtp"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/joho/godotenv"
)

func SendSignupEmail(username, email string) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// googleKey := os.Getenv("GOOGLE_SMTP")

	// auth := smtp.PlainAuth(
	// 	"",
	// 	"jonasmoon23@gmail.com",
	// 	googleKey,
	// 	"smtp.gmail.com",
	// )

	// msg := "Subject: Echo-One-Signup\rCongrats " + username + "! You successfully signed up to Echo-One"

	// err = smtp.SendMail(
	// 	"smtp.gmail.com:587",
	// 	auth,
	// 	"mrmonk@detective.com", // from
	// 	[]string{email},        // slice of addresses to send to
	// 	[]byte(msg),            // message sent as byte array
	// )
	// if err != nil {
	// 	fmt.Println(err)
	// }

	from := mail.NewEmail(os.Getenv("SENDGRID_FROM_USERNAME"), os.Getenv("SENDGRID_FROM_EMAIL"))
	subject := "Sending with Twilio SendGrid is Fun"
	to := mail.NewEmail(username, email)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
