package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type CustomClaims struct {
	Username string `json:"name"`
	Email    string
	jwt.StandardClaims
}

func main() {
	e := echo.New()

	e.POST("/signup", SignUp)
	e.POST("/login", Login)
	// e.POST("/sendmail", SendMailHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start("localhost:1323"))
}
