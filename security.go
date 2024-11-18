package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
)

func enableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "chrome-extension://mkjhflenhpjedegkhgnjlconogccecmp")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == http.MethodOptions {
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

func (cfg *apiConfig) hanlderValidateToken(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Authorization header missing")
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	_, err := validateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		log.Printf("Something went wrong while validating the token: %v\n", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func validateJWT(tokenString string, jwtSecret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
