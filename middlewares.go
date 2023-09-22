package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func withAuth(handler http.HandlerFunc) http.HandlerFunc {
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := bearerToken[1]
		claims := &jwt.StandardClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				log.Println(token)
				return
			}
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		log.Printf("%q %q", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	}
}
