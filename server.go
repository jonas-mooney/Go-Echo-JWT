package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	SignUp(w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	Login(w, r)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signup", signUpHandler)
	r.HandleFunc("/login", loginHandler)
	r.Use(JWT_auth_middleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
