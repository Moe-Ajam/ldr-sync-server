package auth

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func GetClaims(w http.ResponseWriter, r *http.Request, claims jwt.Claims, jwtSecret string) error {
	c, err := r.Cookie("token")
	if err != nil {
		return err
	}

	tknString := c.Value

	token, err := jwt.ParseWithClaims(tknString, claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Token is invalid")
	}
	return nil
}
