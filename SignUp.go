package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"echo-one/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
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
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
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
	} else if userCheck.Username == username || userCheck.Email == email {
		w.WriteHeader(400)
		w.Write([]byte("Username or email unavailable"))
		return
	}

	_, err = db.Exec("INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)", uuid, username, email, hashedString)
	if err != nil {
		fmt.Println("Error occurred:", err)
		w.WriteHeader(500)
		w.Write([]byte("Error creating account"))
	}

	SendSignupEmail(username, email)
	nameTokenJSON, err := CreateJWT(username)
	if err != nil {
		fmt.Println("Error occurred:", err)
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(nameTokenJSON)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}
