package main

import (
	"fmt"
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

	fromUsername := os.Getenv("SENDGRID_FROM_USERNAME")
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
	apiKey := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail(fromUsername, fromEmail)
	subject := "Echo-One-Signup"
	to := mail.NewEmail(username, email)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<h1>Congrats! Welcome to Echo One!</h1>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode, "Email sent successfully")
	}
}
