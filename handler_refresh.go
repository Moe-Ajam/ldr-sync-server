package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Moe-Ajam/ldr-sync-server/internal/auth"
	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Something went wrong!!", err)
		return
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		helpers.RespondWithError(w, http.StatusBadRequest, "Token is still valid, please try again later")
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(cfg.jwtSecret)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	helpers.RespondWithJSON(w, http.StatusOK, tokenString)
}
