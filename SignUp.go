package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"echo-one/models"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func SignUp(c echo.Context) error {
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

	uuid := uuid.New()
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Hashing error")
	}
	hashedString := string(hashedPassword)

	var userCheck models.User

	err = db.QueryRow("SELECT username, email FROM users WHERE username = $1 OR email = $2", username, email).Scan(&userCheck.Username, &userCheck.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Unique username & email passed")
		} else {
			fmt.Println("Error occurred:", err)
		}
	} else if userCheck.Username == username {
		return c.JSON(http.StatusConflict, "That username is taken")
	} else if userCheck.Email == email {
		return c.JSON(http.StatusConflict, "That email is taken")
	}

	_, err = db.Exec("INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)", uuid, username, email, hashedString)
	if err != nil {
		fmt.Println("Error occurred:", err)
	}

	SendSignupEmail(username, email)

	claims := CustomClaims{
		username,
		email,
		jwt.StandardClaims{
			ExpiresAt: jwt.NewTime(15000),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	keyFromEnv := os.Getenv("JWT_SIGNING_KEY")
	mySigningKey := []byte(keyFromEnv)

	ss, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"name":  username,
		"token": ss,
	})
}
