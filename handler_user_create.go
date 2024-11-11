package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Moe-Ajam/ldr-sync-server/internal/database"
	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type registerResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w, r)
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println(err)
		helpers.RespondWithError(w, 500, "something went wrong while creating a user")
		return
	}

	retrievedUser, err := cfg.DB.GetUserByEmail(context.Background(), params.Email)
	if err == nil {
		fmt.Printf("email already exists and has the id: %s\n", retrievedUser.ID)
		helpers.RespondWithError(w, http.StatusConflict, "email already exists")
		return
	}
	retrievedUser, err = cfg.DB.GetUserByName(context.Background(), params.Username)
	if err == nil {
		fmt.Printf("name already exists and has the id: %s\n", retrievedUser.ID)
		helpers.RespondWithError(w, http.StatusConflict, "name already exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		helpers.RespondWithError(w, 500, "something went wrong while creating a user")
		return
	}

	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Username,
		Email:     params.Email,
		Password:  string(hash),
	})

	response := registerResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	helpers.RespondWithJSON(w, http.StatusCreated, response)
}
