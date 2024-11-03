package main

import (
	"fmt"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: Make the token validation into a function in the auth package
func (cfg *apiConfig) handleWelcome(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			helpers.RespondWithError(w, http.StatusUnauthorized, "No cookie present")
			return
		}
		fmt.Printf("Something went wrong while hanlding the request: %v\n", err)
		helpers.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	tknString := c.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tknString, claims, func(token *jwt.Token) (any, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println(tknString)
			helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		fmt.Printf("Something went wrong: %v\n", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if !token.Valid {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Inavalid Token")
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s, with email %s!", claims.Username, claims.Email)))
}
