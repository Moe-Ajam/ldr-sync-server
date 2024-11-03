package main

import (
	"fmt"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/auth"
	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handleWelcome(w http.ResponseWriter, r *http.Request) {
	claims := &Claims{}
	err := auth.GetClaims(w, r, claims, cfg.jwtSecret)
	if err != nil {
		if err == http.ErrNoCookie {
			helpers.RespondWithError(w, http.StatusUnauthorized, "No cookie present")
			return
		}
		if err == jwt.ErrSignatureInvalid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		helpers.RespondWithError(w, http.StatusUnauthorized, "Token is invalid")
		fmt.Println("Something went wrong!!")
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s, with email %s!", claims.Username, claims.Email)))
}
