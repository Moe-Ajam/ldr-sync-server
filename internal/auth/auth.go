package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
)

func GetClaims(w http.ResponseWriter, r *http.Request, claims jwt.Claims, jwtSecret string) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "authorization header missing")
		return errors.New("authorization header missing")
	}

	tknString := strings.TrimPrefix(authHeader, "Bearer ")

	_, err := jwt.ParseWithClaims(tknString, claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		if err == http.ErrNoCookie {
			helpers.RespondWithError(w, http.StatusUnauthorized, "No cookie present")
			return err
		}
		if err == jwt.ErrSignatureInvalid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return err
		}
		helpers.RespondWithError(w, http.StatusUnauthorized, "Token is invalid")
		return err
	}
	return nil
}
