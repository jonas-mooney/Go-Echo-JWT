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

func Login(w http.ResponseWriter, r *http.Request) error {
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
		if err == sql.ErrNoRows {
			w.WriteHeader(401)
			w.Write([]byte("No account with this username and email"))
			return nil
		} else if user.Username == username {
			fmt.Println("Username matches database")
		} else if user.Email == email {
			fmt.Println("Email matches database")
		}
	}

	passwordFromDB := user.Password
	err = bcrypt.CompareHashAndPassword([]byte(passwordFromDB), []byte(password))
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("Incorrect Password"))
		return nil
	}

	nameTokenJSON, err := CreateJWT(user.Username)
	if err != nil {
		fmt.Println("Error occurred:", err)
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(nameTokenJSON)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return err
	}

	return err
}
