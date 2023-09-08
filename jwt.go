package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type ResponseData struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func CreateJWT(username string) ([]byte, error) {
	jwt_key := os.Getenv("JWT_SIGNING_KEY")
	keyByte := []byte(jwt_key)

	claims := CustomClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(keyByte)
	if err != nil {
		fmt.Println("error occured in jwt.go: ", err)
	}

	responseData := ResponseData{
		Username: username,
		Token:    ss,
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func JWT_auth_middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("JWT auth middleware function ran")
		next.ServeHTTP(w, r)
	})
}
