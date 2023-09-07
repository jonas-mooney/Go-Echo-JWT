package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"
	"echo-one/models"

	"github.com/joho/godotenv"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	var user models.User

	err = db.QueryRow("SELECT username, email, password FROM users WHERE username = $1 OR email = $2", username, email).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		fmt.Println("Error occurred:", err)
		w.Write([]byte("An error occured"))
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

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("Login function ran"))
}
