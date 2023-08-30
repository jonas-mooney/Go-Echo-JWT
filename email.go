package main

import (
	"fmt"
	// "net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	// "github.com/labstack/echo/v4"
)

// func SendMailHandler(c echo.Context) error {
// 	Send()
// 	return c.JSON(http.StatusOK, echo.Map{
// 		"Mail Status": "Email sent",
// 	})
// }

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
