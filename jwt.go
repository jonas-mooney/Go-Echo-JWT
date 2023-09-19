package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

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

func JWT_auth(w http.ResponseWriter, r *http.Request) error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	jwt_key := os.Getenv("JWT_SIGNING_KEY")

	tokenString := r.Header.Get("Token")
	if tokenString == "" {
		w.WriteHeader(401)
		w.Write([]byte("No jwt present"))
		return err
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if err != nil {
			return err, err // update
		}
		return []byte(jwt_key), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		w.WriteHeader(200)
		w.Write([]byte(claims.Username))
		return nil
	} else {
		w.WriteHeader(401)
		return err
	}
}

// test case where the user sends a token that contains a string not resembling a jwt
