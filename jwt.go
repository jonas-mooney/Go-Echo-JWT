package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	// "github.com/joho/godotenv"
)

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type ValidStatus struct {
	Username      string `json:"username"`
	Authenticated bool   `json:"authenticated"`
}

type Config struct {
	JWTSigningKey string
}

var cfg Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg.JWTSigningKey = os.Getenv("JWT_SIGNING_KEY")
}

func CreateJWT(username string) ([]byte, error) {
	jwt_key := cfg.JWTSigningKey
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
		log.Printf("Error signing token: %v", err)
	}

	responseData := TokenResponse{
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
	jwt_key := cfg.JWTSigningKey

	tokenString := r.Header.Get("Token")
	if tokenString == "" {
		w.WriteHeader(401)
		w.Write([]byte("No jwt present"))
		return nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwt_key), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return nil
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {

		responseData := ValidStatus{
			Username:      claims.Username,
			Authenticated: true,
		}

		jsonData, err := json.Marshal(responseData)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(jsonData))
		return nil
	} else {
		w.WriteHeader(401)
		return err
	}
}
