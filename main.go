package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	err := JWT_auth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	err := SignUp(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := Login(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth", authHandler)
	r.HandleFunc("/signup", signUpHandler)
	r.HandleFunc("/login", loginHandler)

	log.Printf("Running on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
	// log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
