package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	// "github.com/joho/godotenv"

	"echo-one/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func SignUp(w http.ResponseWriter, r *http.Request) error {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Printf("Error loading .env file: %v", err)
	// }

	connStr := os.Getenv("RAILWAY_PG_CONNECTION_STRING")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	if connStr == "" {
		w.WriteHeader(500)
		w.Write([]byte("RAILWAY_PG_CONNECTION_STRING environment variable not set"))
		return nil
	}

	uuid := uuid.New()
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
	}
	hashedString := string(hashedPassword)

	var userCheck models.User

	err = db.QueryRow("SELECT username, email FROM users WHERE username = $1 OR email = $2", username, email).Scan(&userCheck.Username, &userCheck.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Unique username & email passed")

			/// just a test
			// w.WriteHeader(200)
			// w.Write([]byte("Success connecting and querying db"))
			// return err
			///

		} else {
			log.Printf("Error occurred: %v", err)
		}
	} else if userCheck.Username == username || userCheck.Email == email {
		w.WriteHeader(400)
		w.Write([]byte("Username or email unavailable"))
		return err
	}

	_, err = db.Exec("INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)", uuid, username, email, hashedString)
	if err != nil {
		log.Printf("Error occurred in signup: %v", err)
		w.WriteHeader(500)
		w.Write([]byte("Error creating account from signup.go 74"))
		return err
	}

	err = SendSignupEmail(username, email)
	if err != nil {
		return err
	}

	nameTokenJSON, err := CreateJWT(username)
	if err != nil {
		log.Printf("Error occurred: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(nameTokenJSON)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return err
	}

	return err
}
