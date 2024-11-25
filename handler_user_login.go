package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// TODO: Make the token creation into its own function in an auth package
func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
	// if r.Method == http.MethodOptions {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	log.Printf("Login params: %v\n", params)
	if err != nil {
		fmt.Printf("something went wrong while decoding the params for the login: %v\n", err)
		helpers.RespondWithError(w, 500, "something went wrong, could not login")
		return
	}

	user, err := cfg.DB.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		fmt.Printf("user with the email %s does not exist\n", params.Email)
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unautorized")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		fmt.Printf("the password %s is wrong for the email: %s\n", params.Password, params.Email)
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unautorized")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Name,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		fmt.Printf("there was a problem signing the token: %v\n", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "something went wrong, could not create user")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		SameSite: http.SameSiteNoneMode,
	})

	response := LoginResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: tokenString,
	}

	helpers.RespondWithJSON(w, http.StatusOK, response)
}
