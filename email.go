package main

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/joho/godotenv"
)

type Configuration struct {
	SendGridFromUsername string
	SendGridFromEmail    string
	SendGridAPIKey       string
}

func loadConfig() (Configuration, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg := Configuration{
		SendGridFromUsername: os.Getenv("SENDGRID_FROM_USERNAME"),
		SendGridFromEmail:    os.Getenv("SENDGRID_FROM_EMAIL"),
		SendGridAPIKey:       os.Getenv("SENDGRID_API_KEY"),
	}

	if cfg.SendGridFromUsername == "" || cfg.SendGridFromEmail == "" || cfg.SendGridAPIKey == "" {
		return Configuration{}, err
	}

	return cfg, nil
}

func SendSignupEmail(username, email string) {
	cfg, err := loadConfig()
	if err != nil {
		log.Println("Error loading config: ", err)
	}

	from := mail.NewEmail(cfg.SendGridFromUsername, cfg.SendGridFromEmail)
	subject := "Echo-One-Signup"
	to := mail.NewEmail(username, email)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<h1>Congrats " + username + "! " + "Welcome to Echo One!</h1>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(cfg.SendGridAPIKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println("Sendgrid signup email error: ", err)
	} else if response.StatusCode != 202 {
		log.Println("Sendgrid signup email status: ", response.StatusCode)
	}
}

/*
Testing Environment:
Testing Sendgrip api key
*/

/*
Production Environment:
In production, you should have monitoring and alerting systems in place to detect and respond to issues with your email sending functionality. This can include:
Health Checks: Implement a health check endpoint in your application that periodically verifies the SendGrid API's operational status. This endpoint can be monitored by external systems or load balancers, which can take actions based on the results (e.g., switch to a backup email provider if SendGrid is down).
Alerting: Set up alerts to notify you when SendGrid experiences issues, such as rate limits exceeded, outages, or significant delivery failures. Services like SendGrid often provide APIs or webhooks to facilitate this kind of monitoring.
Logging and Error Handling: Your application should have comprehensive error handling and logging for email sending. When errors occur, they should be logged and possibly reported so that you can investigate and address the issue promptly.
Fallback Mechanism: Consider having a fallback mechanism in place. If SendGrid experiences issues, you can switch to an alternative email provider or service to ensure that user signups are not interrupted.
*/
