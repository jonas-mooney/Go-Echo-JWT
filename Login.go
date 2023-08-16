package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"
	"echo-one/models"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func Login(c echo.Context) error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	connStr := os.Getenv("RAILWAY_PG_CONNECTION_STRING")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// token := c.Request().Header.Get("Authorization")
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User

	err = db.QueryRow("SELECT username, email, password FROM users WHERE username = $1 OR email = $2", username, email).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		fmt.Println("Error occurred:", err)
		return c.JSON(http.StatusInternalServerError, "An error occurred")
	}

	if err == sql.ErrNoRows {
		fmt.Println("Unique username & email passed")
	} else if user.Username == username {
		fmt.Println("Username matches database")
	} else if user.Email == email {
		fmt.Println("Email matches database")
	}

	passwordFromDB := user.Password
	err = bcrypt.CompareHashAndPassword([]byte(passwordFromDB), []byte(password))
	if err == nil {
		fmt.Println("Password matches")
	} else {
		fmt.Println("Password doesn't match")
	}

	claims := CustomClaims{
		username,
		email,
		jwt.StandardClaims{
			ExpiresAt: jwt.NewTime(15000),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	mySigningKey := []byte("h6t5rd3s4a12h")

	ss, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		println(err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"name":  username,
		"token": ss,
	})

}
