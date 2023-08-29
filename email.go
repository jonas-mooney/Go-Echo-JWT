package main

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
)

func SendMailHandler(c echo.Context) error {
	Send()
	return c.JSON(http.StatusOK, echo.Map{
		"Mail Status": "Email sent",
	})
}

func Send() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	googleKey := os.Getenv("GOOGLE_SMTP")
	auth := smtp.PlainAuth(
		"",
		"placeholder@gma.com",
		googleKey,
		"smtp.gmail.com",
	)

	msg := "Subject: Cheese recipes\rEmail body also about cheese recipes"

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"placeholder@gma.com",           // from
		[]string{"placeholder@gma.com"}, // slice of addresses to send to
		[]byte(msg),                     // message sent as byte array
	)
	if err != nil {
		fmt.Println(err)
	}
}
