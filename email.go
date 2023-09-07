package main

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendSignupEmail(username, email string) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	googleKey := os.Getenv("GOOGLE_SMTP")

	auth := smtp.PlainAuth(
		"",
		"jonasmooney2@gmail.com",
		googleKey,
		"smtp.gmail.com",
	)

	msg := "Subject: Echo-One-Signup\rCongrats " + username + "! You successfully signed up to Echo-One"

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"jonasmooney2@gmail.com", // from
		[]string{email},          // slice of addresses to send to
		[]byte(msg),              // message sent as byte array
	)
	if err != nil {
		fmt.Println(err)
	}
}
