package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

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
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/signup", signUpHandler)
	r.HandleFunc("/login", loginHandler)
	// r.Use(JWT_auth_middleware)

	log.Printf("Running on :8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
